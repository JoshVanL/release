/*
Copyright 2021 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registry

import (
	"context"
	"log"

	"github.com/cert-manager/release/pkg/release/docker"
	"github.com/cert-manager/release/pkg/release/images"
)

func CreateManifestList(ctx context.Context, name string, tars []images.Tar) error {
	imageNames := make([]string, len(tars))
	for i, t := range tars {
		imageNames[i] = t.ImageName()
	}

	log.Printf("Creating manifest list %q", name)
	if err := docker.CreateManifestList(ctx, name, imageNames); err != nil {
		log.Printf("Failed to create manifest list with name %q - ensure no existing manifest list exists with the same name, and ensure all member images are pushed to the remote registry.", name)
		return err
	}

	for _, t := range tars {
		a := manifestListAnnotationsForOSArch(t.OS(), t.Architecture())
		log.Printf("Annotating image %q with os=%q, arch=%q, variant=%q", t.ImageName(), a.os, a.arch, a.variant)
		if err := docker.AnnotateManifestList(ctx, name, t.ImageName(), a.os, a.arch, a.variant); err != nil {
			log.Printf("Failed to annotate manifest list with os/arch information.")
			return err
		}
	}

	log.Printf("Created manifest list %q", name)
	return nil
}

type manifestAnnotation struct {
	os, arch, variant string
}

func manifestListAnnotationsForOSArch(os, arch string) manifestAnnotation {
	if arch == "arm" {
		return manifestAnnotation{
			os:      os,
			arch:    arch,
			variant: "v7",
		}
	}
	if arch == "arm64" {
		return manifestAnnotation{
			os:      os,
			arch:    arch,
			variant: "v8",
		}
	}
	return manifestAnnotation{
		os:   os,
		arch: arch,
	}
}
