package main

import (
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
	InValidImage := strings.Contains(Image, "pivotal.io")
	DevImage := strings.Contains(Image, "dev")
	if InValidImage || !DevImage {
		log.Fatalln("Please provide image reference which points to dev.registry.tanzu.vmware.com")
	} else {
		log.Println("Image reference in package CR is validated successfully")
	}
}
