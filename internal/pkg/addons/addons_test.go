// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package addons

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/aws/amazon-ecs-cli-v2/internal/pkg/addons/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddons_Template(t *testing.T) {
	const testSvcName = "mysvc"
	testErr := errors.New("some error")
	testCases := map[string]struct {
		mockAddons func(ctrl *gomock.Controller) *Addons

		wantedTemplate string
		wantedErr      error
	}{
		"return ErrDirNotExist if addons doesn't exist in a service": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).
					Return(nil, testErr)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedErr: &ErrDirNotExist{
				SvcName:   testSvcName,
				ParentErr: testErr,
			},
		},
		"return err on invalid Metadata fields": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).Return([]string{"first.yaml", "invalid-metadata.yaml"}, nil)

				first, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "first.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "first.yaml").Return(first, nil)

				second, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "invalid-metadata.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "invalid-metadata.yaml").Return(second, nil)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedErr: errors.New(`merge addon invalid-metadata.yaml under service mysvc: metadata key "Services" already exists with a different definition`),
		},
		"returns err on invalid Parameters fields": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).Return([]string{"first.yaml", "invalid-parameters.yaml"}, nil)

				first, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "first.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "first.yaml").Return(first, nil)

				second, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "invalid-parameters.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "invalid-parameters.yaml").Return(second, nil)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedErr: errors.New(`merge addon invalid-parameters.yaml under service mysvc: parameter logical ID "Name" already exists with a different definition`),
		},
		"returns err on invalid Mappings fields": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).Return([]string{"first.yaml", "invalid-mappings.yaml"}, nil)

				first, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "first.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "first.yaml").Return(first, nil)

				second, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "invalid-mappings.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "invalid-mappings.yaml").Return(second, nil)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedErr: errors.New(`merge addon invalid-mappings.yaml under service mysvc: mapping "MyTableDynamoDBSettings.test" already exists with a different definition`),
		},
		"returns err on invalid Conditions fields": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).Return([]string{"first.yaml", "invalid-conditions.yaml"}, nil)

				first, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "first.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "first.yaml").Return(first, nil)

				second, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "invalid-conditions.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "invalid-conditions.yaml").Return(second, nil)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedErr: errors.New(`merge addon invalid-conditions.yaml under service mysvc: condition "IsProd" already exists with a different definition`),
		},
		"returns err on invalid Resources fields": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).Return([]string{"first.yaml", "invalid-resources.yaml"}, nil)

				first, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "first.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "first.yaml").Return(first, nil)

				second, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "invalid-resources.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "invalid-resources.yaml").Return(second, nil)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedErr: errors.New(`merge addon invalid-resources.yaml under service mysvc: resource "MyTable" already exists with a different definition`),
		},
		"returns err on invalid Outputs fields": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).Return([]string{"first.yaml", "invalid-outputs.yaml"}, nil)

				first, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "first.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "first.yaml").Return(first, nil)

				second, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "invalid-outputs.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "invalid-outputs.yaml").Return(second, nil)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedErr: errors.New(`merge addon invalid-outputs.yaml under service mysvc: output "MyTableAccessPolicy" already exists with a different definition`),
		},
		"merge fields successfully": {
			mockAddons: func(ctrl *gomock.Controller) *Addons {
				ws := mocks.NewMockworkspaceReader(ctrl)
				ws.EXPECT().ReadAddonsDir(testSvcName).Return([]string{"first.yaml", "second.yaml"}, nil)

				first, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "first.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "first.yaml").Return(first, nil)

				second, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "second.yaml"))
				ws.EXPECT().ReadAddon(testSvcName, "second.yaml").Return(second, nil)
				return &Addons{
					svcName: testSvcName,
					ws:      ws,
				}
			},
			wantedTemplate: func() string {
				wanted, _ := ioutil.ReadFile(filepath.Join("testdata", "merge", "wanted.yaml"))
				return string(wanted)
			}(),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// GIVEN
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			addons := tc.mockAddons(ctrl)

			// WHEN
			actualTemplate, actualErr := addons.Template()

			// THEN
			if tc.wantedErr != nil {
				require.EqualError(t, actualErr, tc.wantedErr.Error())
			} else {
				require.NoError(t, actualErr)
				require.Equal(t, tc.wantedTemplate, actualTemplate)
			}
		})
	}
}
