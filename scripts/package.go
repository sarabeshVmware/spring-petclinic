package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gitlab.eng.vmware.com/tap/tap-packages/scripts/pkg"
	"gopkg.in/yaml.v3"
)

type PackageCR struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		RefName    string `yaml:"refName"`
		Version    string `yaml:"version" validate:"required"`
		ReleasedAt string `yaml:"releasedAt" validate:"ISO8601date"`
		Template   struct {
			Spec struct {
				Fetch []struct {
					ImgpkgBundle struct {
						Image string `yaml:"image"`
					} `yaml:"imgpkgBundle"`
				} `yaml:"fetch"`
				Template []struct {
					Kbld struct {
						Paths []string `yaml:"paths"`
					} `yaml:"kbld,omitempty"`
					Ytt struct {
						Paths []string `yaml:"paths"`
					} `yaml:"ytt,omitempty"`
				} `yaml:"template"`
				Deploy []struct {
					Kapp struct {
					} `yaml:"kapp"`
				} `yaml:"deploy"`
			} `yaml:"spec"`
		} `yaml:"template" validate:"required"`
	} `yaml:"spec"`
}

type CraneManifestsOutput struct {
	SchemaVersion int    `yaml:"schemaVersion"`
	MediaType     string `yaml:"mediaType"`
	Config        struct {
		MediaType string `yaml:"mediaType"`
		Size      int    `yaml:"size"`
		Digest    string `yaml:"digest"`
	} `yaml:"config"`
	Layers []struct {
		MediaType string `yaml:"mediaType"`
		Size      int    `yaml:"size"`
		Digest    string `yaml:"digest"`
	} `yaml:"layers"`
	Manifests []struct {
		MediaType string `yaml:"mediaType"`
		Digest    string `yaml:"digest"`
		Size      int    `yaml:"size"`
		Platform  struct {
			Architecture string `yaml:"architecture"`
			Os           string `yaml:"os"`
		} `yaml:"platform"`
	} `yaml:"manifests"`
}

// type ImgPkgDescribeOutput struct {
// 	Images []struct {
// 		Image       string `yaml:"Image"`
// 		Type        string `yaml:"Type"`
// 		Origin      string `yaml:"Origin"`
// 		// Annotations struct {
// 		// 	KbldCarvelDevID string `yaml:"kbld.carvel.dev/id"`
// 		// } `yaml:"Annotations"`
// 		// Tag string `yaml:"tag,omitempty"`
// 		// URL string `yaml:"url,omitempty"`
// 	} `yaml:"Images"`
// }
// type ImgPkgDescribeOutput struct {
// 	Sha     string `yaml:"sha"`
// 	Content struct {
// 		Images []struct {
// 		   sha256 struct {
// 		   	  Image     string `yaml:"image"`
// 		   }`yaml:"sha256"`

// 		} `yaml:"images"`
// 	} `yaml:"content"`
// }
type Imagesha struct {
	// Annotations struct {
	// 	KbldCarvelDevID      string `yaml:"kbld.carvel.dev/id"`
	// 	KbldCarvelDevOrigins string `yaml:"kbld.carvel.dev/origins"`
	// } `yaml:"annotations"`
	Image string `yaml:"image"`
	// ImageType string `yaml:"imageType"`
	// Origin    string `yaml:"origin"`
}
type ImgPkgDescribeOutput struct {
	Sha     string `yaml:"sha"`
	Content struct {
		Images map[string]Imagesha `yaml:"images"`
	} `yaml:"content"`
}

func GetPackageFile(fpath string) PackageCR {
	inputBytes, err := os.ReadFile(fpath)
	pkg.CheckError(err)
	pkgFile := PackageCR{}
	err = yaml.Unmarshal(inputBytes, &pkgFile)
	pkg.CheckError(err)
	return pkgFile
}

var validate *validator.Validate

func IsISO8601Date(fl validator.FieldLevel) bool {
	ISO8601DateRegexString := "^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])(?:T|\\s)(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])?(Z)?$"
	ISO8601DateRegex := regexp.MustCompile(ISO8601DateRegexString)
	return ISO8601DateRegex.MatchString(fl.Field().String())
}

func main() {
	fpath := os.Args[1]
	log.Printf("Validating file extention...")
	pkg.CheckFileExtension(fpath, ".yaml")
	packageFile := GetPackageFile(fpath)
	log.Printf("Validating Image reference in package CR file: %s", fpath)
	ValidateImage(packageFile.Spec.Template.Spec.Fetch[0].ImgpkgBundle.Image)
	log.Printf("Validating template section in package CR file: %s", fpath)
	ValidateTemplateSection(packageFile)
	log.Printf("Validating package CR file: %s", fpath)
	validate = validator.New()
	log.Printf("Registring custom validator for ISO8601 date validation")
	validate.RegisterValidation("ISO8601date", IsISO8601Date)
	err := validate.Struct(packageFile)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Println(validationErrors)
		for _, err := range validationErrors {
			if err.Tag() == "ISO8601date" {
				log.Fatalln("Please provide field ", err.StructNamespace(), "in format", err.Tag())
			} else {
				log.Fatalln("Field ", err.StructNamespace(), "is", err.Tag())
			}
			log.Println()
		}
	} else {
		log.Println("package CR file:", fpath, "validated successfully")
	}
}

func ValidateImage(Image string) {
	log.Println("Validating Image", Image)
	InValidImage := strings.Contains(Image, "pivotal.io")
	DevImage := strings.Contains(Image, "dev.registry.tanzu.vmware.com")
	ShaImage := strings.Contains(Image, "@sha256:")
	ImageIndex, _ := ValidateImageIndex(Image)
	if InValidImage || !DevImage {
		log.Fatalln("Please provide image reference which points to dev.registry.tanzu.vmware.com")
	}
	if !ShaImage {
		log.Fatalln("Please provide bundle image in digested form instead of tag")
	}
	if ImageIndex {
		log.Fatalln("Please provide bundle image which does not refer to multiple image indexes")
	}
	log.Println("Image reference in package CR is validated successfully")
}

func ValidateTemplateSection(packageFile PackageCR) {
	isKbldPresent, isPathPresentInKbldPath := false, false
	for _, template := range packageFile.Spec.Template.Spec.Template {
		if len(template.Kbld.Paths) != 0 {
			isKbldPresent = true
			for _, path := range template.Kbld.Paths {
				if path == ".imgpkg/images.yml" {
					isPathPresentInKbldPath = true
					break
				}
			}
		}
	}
	if !isKbldPresent {
		log.Fatal("kbld entry is absent in the template section.")
	}
	if !isPathPresentInKbldPath {
		log.Fatal(`".imgpkg/images.yml" entry is absent in the kbld, in the template section.`)
	}
	log.Print("template section validated successfully.")
}

func ValidateImageIndex(Image string) (bool, error) {
	imgpkgDescribeInfo := ImgPkgDescribeOutput{}
	log.Println("Executing imgpkg describe command...")
	cmd := fmt.Sprintf("imgpkg describe -b %s -o yaml", Image)
	output, err := pkg.ExecuteCmd(cmd)
	status := false
	if err != nil {
		log.Printf("error while running imgpkg describe %s ", Image)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
		return status, err
	} 
	resb := []byte(output)
	err = yaml.Unmarshal(resb, &imgpkgDescribeInfo)
	if err != nil {
		log.Printf("error while unmarshal imgpkg describe output %s ", Image)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
		return status, err
	}
	var multiple []string
	for _,v := range(imgpkgDescribeInfo.Content.Images){
		craneManifestsInfo := CraneManifestsOutput{}
		cmd := fmt.Sprintf("crane manifest %s", v.Image)
		output, err = pkg.ExecuteCmd(cmd)
		if err != nil {
			log.Printf("error while running crane manifest %s ", v.Image)
			log.Printf("error: %s", err)
			log.Printf("output: %s", output)
		} 
		resb2 := []byte(output)
		err = yaml.Unmarshal(resb2, &craneManifestsInfo)
		if err != nil {
			log.Printf("error while unmarshal crane manifest output %s ", v.Image)
			log.Printf("error: %s", err)
			log.Printf("output: %s", output)
			return status, err
		}
		if len(craneManifestsInfo.Manifests) > 0{
			for _, manifests := range(craneManifestsInfo.Manifests){
				log.Printf("Manifest Digest: %s", manifests.Digest)
				log.Printf("Manifest Platform Architecture: %s", manifests.Platform.Architecture)
				log.Printf("Manifest Platform Os: %s", manifests.Platform.Os)
				log.Printf("Imgpkg bundle %s, image %s contains multiple image index which is not supported by Tanzunet.", Image, v.Image)
				status = true
			}
			multiple = append(multiple, v.Image)
		} else {
			log.Printf("Manifests list is not found for image %s. Image does not contain any multiple image indexs", v.Image)
		}
	}
	if status{
		log.Printf("Multiple image index found for images : %+v\n", multiple)
	}
	
	return status, err
}
