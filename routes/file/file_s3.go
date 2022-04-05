/**
* This file is part of VILLASweb-backend-go
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
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/request"
	"io"
	"net/url"
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

func getS3Session() (*session.Session, string, error) {

	bucket, err := configuration.GlobalConfig.String("s3.bucket")
	if err != nil || bucket == "" {
		return nil, "", fmt.Errorf("no S3 bucket configured: %s", err)
	}

	if s3Session == nil {
		var err error
		s3Session, err = createS3Session()
		if err != nil {
			return nil, "", err
		}
	}

	return s3Session, bucket, nil
}

func createS3Session() (*session.Session, error) {
	endpoint, err := configuration.GlobalConfig.String("s3.endpoint")
	if err != nil {
		return nil, err
	}
	region, err := configuration.GlobalConfig.String("s3.region")
	if err != nil {
		return nil, err
	}
	pathStyle, err := configuration.GlobalConfig.Bool("s3.pathstyle")
	if err != nil {
		return nil, err
	}
	nossl, err := configuration.GlobalConfig.Bool("s3.nossl")
	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(
		&aws.Config{
			Region:           aws.String(region),
			Endpoint:         aws.String(endpoint),
			DisableSSL:       aws.Bool(nossl),
			S3ForcePathStyle: aws.Bool(pathStyle),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return sess, nil
}

func (f *File) putS3(fileContent io.Reader) error {

	// The session the S3 Uploader will use
	sess, bucket, err := getS3Session()
	if err != nil {
		return err
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
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

func (f *File) getS3Url() (string, error) {

	// The session the S3 Uploader will use
	sess, bucket, err := getS3Session()
	if err != nil {
		return "", err
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

	err = updateS3Request(req)
	if err != nil {
		return "", err
	}

	urlStr, err := req.Presign(5 * 24 * 60 * time.Minute)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}

//lint:ignore U1000 will be used later
func (f *File) deleteS3() error {

	// The session the S3 Uploader will use
	sess, bucket, err := getS3Session()
	if err != nil {
		return err
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

// updateS3Request updates the request host to the public accessible S3
// endpoint host so that presigned URLs are still valid when accessed
// by the user
func updateS3Request(req *request.Request) error {
	epURL, err := getS3EndpointURL()
	if err != nil {
		return err
	}

	req.HTTPRequest.URL.Scheme = epURL.Scheme
	req.HTTPRequest.URL.Host = epURL.Host

	return nil
}

func getS3EndpointURL() (*url.URL, error) {
	ep, err := configuration.GlobalConfig.String("s3.endpoint-public")
	if err != nil {
		ep, err = configuration.GlobalConfig.String("s3.endpoint")
		if err != nil {
			return nil, errors.New("missing s3.endpoint setting")
		}
	}

	epURL, err := url.Parse(ep)
	if err != nil {
		return nil, err
	}

	return epURL, nil
}
