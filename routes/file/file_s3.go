/** File package, S3 uploads.
*
* @author Steffen Vogel <svogel2@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package file

import (
	"fmt"
	"io"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

// Global session
var s3Session *session.Session = nil

func getS3Session() (*session.Session, error) {
	if s3Session == nil {
		var err error
		s3Session, err = createS3Session()
		if err != nil {
			return nil, err
		}
	}

	return s3Session, nil
}

func createS3Session() (*session.Session, error) {
	endpoint, err := configuration.GolbalConfig.String("s3.endpoint")
	region, err := configuration.GolbalConfig.String("s3.region")
	pathStyle, err := configuration.GolbalConfig.Bool("s3.pathstyle")
	nossl, err := configuration.GolbalConfig.Bool("s3.nossl")

	sess, err := session.NewSession(
		&aws.Config{
			Region:           aws.String(region),
			Endpoint:         aws.String(endpoint),
			DisableSSL:       aws.Bool(nossl),
			S3ForcePathStyle: aws.Bool(pathStyle),
		},
	)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (f *File) putS3(fileContent io.Reader) error {

	bucket, err := configuration.GolbalConfig.String("s3.bucket")
	if err != nil || bucket == "" {
		return fmt.Errorf("No S3 bucket configured")
	}

	// The session the S3 Uploader will use
	sess, err := getS3Session()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f.Key = uuid.New().String()
	f.FileData = nil

	// Upload the file to S3.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(f.Key),
		Body:   fileContent,
	})
	if err != nil {
		return fmt.Errorf("Failed to upload file, %v", err)
	}

	return nil
}

func (f *File) getS3Url() (string, error) {
	bucket, err := configuration.GolbalConfig.String("s3.bucket")
	if err != nil || bucket == "" {
		return "", fmt.Errorf("No S3 bucket configured")
	}

	// The session the S3 Uploader will use
	sess, err := getS3Session()
	if err != nil {
		return "", fmt.Errorf("Failed to create session: %s", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(f.Key),
		ResponseContentType:        aws.String(f.Type),
		ResponseContentDisposition: aws.String("attachment; filename=" + f.Name),
		// ResponseContentEncoding: aws.String(),
		// ResponseContentLanguage: aws.String(),
		// ResponseCacheControl:    aws.String(),
		// ResponseExpires:         aws.String(),
	})

	urlStr, err := req.Presign(5 * 24 * 60 * time.Minute)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}

func (f *File) deleteS3() error {
	bucket, err := configuration.GolbalConfig.String("s3.bucket")
	if err != nil || bucket == "" {
		return fmt.Errorf("No S3 bucket configured")
	}

	// The session the S3 Uploader will use
	sess, err := getS3Session()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}

	// Create S3 service client
	svc := s3.New(sess)

	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(f.Key),
	})
	if err != nil {
		return err
	}

	f.Key = ""

	return nil
}
