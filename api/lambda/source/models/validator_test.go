package models

/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/require"
)

func TestValidateIntegrationLabelSucceeds(t *testing.T) {
	validator, err := Validator()
	require.NoError(t, err)
	err = validator.Struct(&GetIntegrationTemplateInput{
		AWSAccountID:     aws.String("123456789012"),
		IntegrationLabel: aws.String("Test12- "),
		IntegrationType:  aws.String(IntegrationTypeAWS3),
	})
	require.NoError(t, err)
}

func TestValidateIntegrationLabelFails(t *testing.T) {
	validator, err := Validator()
	require.NoError(t, err)
	err = validator.Struct(&GetIntegrationTemplateInput{
		AWSAccountID:     aws.String("123456789012"),
		IntegrationLabel: aws.String(" "),
		IntegrationType:  aws.String(IntegrationTypeAWS3),
	})
	errorMsg := "Key: 'GetIntegrationTemplateInput.IntegrationLabel'" +
		" Error:Field validation for 'IntegrationLabel' failed on the 'integrationLabel' tag"
	require.EqualError(t, err, errorMsg)
}

func TestValidateNotKmsKey(t *testing.T) {
	validator, err := Validator()
	require.NoError(t, err)
	err = validator.Struct(&PutIntegrationInput{
		PutIntegrationSettings: PutIntegrationSettings{
			AWSAccountID:     aws.String("123456789012"),
			IntegrationLabel: aws.String("Test12- "),
			IntegrationType:  aws.String(IntegrationTypeAWS3),
			UserID:           aws.String("cb7663c7-80ed-420b-a287-ed7dc50a0bf7"),
			KmsKey:           aws.String("not-a-key"),
		},
	})

	errorMsg := "Key: 'PutIntegrationInput.PutIntegrationSettings.KmsKey' " +
		"Error:Field validation for 'KmsKey' failed on the 'kmsKeyArn' tag"
	require.EqualError(t, err, errorMsg)
}

func TestValidateKmsKey(t *testing.T) {
	validator, err := Validator()
	require.NoError(t, err)
	err = validator.Struct(&PutIntegrationInput{
		PutIntegrationSettings: PutIntegrationSettings{
			AWSAccountID:     aws.String("123456789012"),
			IntegrationLabel: aws.String("Test12- "),
			IntegrationType:  aws.String(IntegrationTypeAWS3),
			UserID:           aws.String("cb7663c7-80ed-420b-a287-ed7dc50a0bf7"),
			KmsKey:           aws.String("arn:aws:kms:eu-west-1:415773754570:key/7abf9aaf-0228-4c09-ae6c-c9a0c65e4894"),
		},
	})
	require.NoError(t, err)
}
