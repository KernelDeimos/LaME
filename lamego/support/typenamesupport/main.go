package typenamesupport

import (
	"fmt"
	"github.com/KernelDeimos/LaME/lamego/target"
	"strings"
)

type TypeInformation struct {
	// IsPrimitive is equal to t.TypeOfType == target.PrimitiveType
	IsPrimitive bool

	// IsSpecialVoid is true if the type had no identifier, and thus void
	// type or error should be assumed
	IsSpecialVoid bool

	// FailedToMatch is true if the type is not a primitive but no match
	// was found in the provided mapj
	FailedToMatch bool

	// TypeName is the past token in the dotted type path
	TypeName string

	// IsCurrentPackage is true if the dotted type path has only one
	// token, and the type is not a primitive
	IsCurrentPackage bool

	// LaMEPackage is the dotted string before the type name,
	// without the trailing dot
	LaMEPackage string

	// LaMEPackageMatched is the portion of the dotted string that was
	// matched in the provided map
	LaMEPackageMatched string

	// LanguagePackage is a package prefix recognized from the given
	// packages map, or void
	LanguagePackage string

	// LanguagePackageRemainder is the tokens of the dotted package path
	// that remained after matching the LaME package with a package
	// specified in the packages map
	LanguageRemainder []string
}

func GetTypeInfo(
	// Type to get information for
	t target.Type,
	// Map of LaME packages to a package name or identifier
	// which the user of this function understands
	packages map[string]string,
	// Dotted path of a valid LaME package to assume context
	// for returned values such as IsCurrentPackage
	referencePackage string,
) TypeInformation {
	ret := TypeInformation{}

	ret.IsPrimitive = t.TypeOfType == target.PrimitiveType
	if ret.IsPrimitive {
		ret.TypeName = t.Identifier
		return ret
	}

	parts := strings.Split(t.Identifier, ".")

	// No parts means empty type name; this is "SpecialVoid"
	if len(parts) < 1 {
		ret.IsPrimitive = true
		ret.IsSpecialVoid = true
		return ret
	}

	// One token assumes a type in the current package
	if len(parts) == 1 {
		ret.IsCurrentPackage = true
		ret.TypeName = parts[0]
		return ret
	}

	// parts, but without type name
	pkgParts := parts[:len(parts)-1]
	ret.LaMEPackage = strings.Join(pkgParts, ".")
	ret.TypeName = parts[len(parts)-1]

	// Skip steps below if path matches current package
	fmt.Println(ret.LaMEPackage, referencePackage)
	if ret.LaMEPackage == referencePackage {
		ret.IsCurrentPackage = true
		return ret
	}

	// Begin searching for package associations with provided map
	for rootPkg, targetPkg := range packages {
		if !strings.HasPrefix(ret.LaMEPackage, rootPkg) {
			continue
		}
		rootPkgParts := strings.Split(rootPkg, ".")
		remainderPkgParts := pkgParts[len(rootPkgParts):]
		ret.LanguagePackage = targetPkg
		ret.LanguageRemainder = remainderPkgParts
		ret.LaMEPackageMatched =
			strings.Join(pkgParts[:len(rootPkgParts)], ".")
		return ret
	}

	ret.FailedToMatch = true
	return ret
}
