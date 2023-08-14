package manager

import "github.com/ed3899/kumo/common/interfaces"

func ManagerPathExecutableIsPresentWith(
	utilsFileIsFilePresent func(string) bool,
) ManagerPathExecutableIsPresent {
	managerPathExecutableIsPresent := func(manager interfaces.IClone[Manager]) bool {
		managerClone := manager.Clone()
		return utilsFileIsFilePresent(managerClone.Path().Executable())
	}

	return managerPathExecutableIsPresent
}

type ManagerPathExecutableIsPresent func(interfaces.IClone[Manager]) bool