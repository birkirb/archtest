package archtest_test

import (
	"github.com/mattmcnew/archtest"
	"strings"
	"testing"
)

func TestPackage_ShouldNotDependOn(t *testing.T) {

	t.Run("Succeeds on non dependencies", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/nodependency")

		assertNoError(t, mockT)
	})

	t.Run("Fails on dependencies", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/dependency")

		assertError(t, mockT,
			"github.com/mattmcnew/archtest/examples/testpackage",
			"github.com/mattmcnew/archtest/examples/dependency")
	})

	t.Run("Fails on transative dependencies", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/transative")

		assertError(t, mockT,
			"github.com/mattmcnew/archtest/examples/testpackage",
			"github.com/mattmcnew/archtest/examples/dependency",
			"github.com/mattmcnew/archtest/examples/transative")
	})

	t.Run("Supports multiple packages at once", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/dontdependonanything", "github.com/mattmcnew/archtest/examples/testpackage").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/dependency")

		assertError(t, mockT,
			"github.com/mattmcnew/archtest/examples/testpackage",
			"github.com/mattmcnew/archtest/examples/dependency")
	})

	t.Run("Supports wildcard matching", func(t *testing.T) {
		mockT := new(testingT)
		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/...").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/nodependency")

		assertNoError(t, mockT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage/...").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/nesteddependency")

		assertError(t, mockT, "github.com/mattmcnew/archtest/examples/testpackage/nested/dep", "github.com/mattmcnew/archtest/examples/nesteddependency")
	})

	t.Run("Supports checking imports in test files", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage/...").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/testfiledeps/testonlydependency")

		assertNoError(t, mockT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage/...").
			IncludeTests().
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/testfiledeps/testonlydependency")

		assertError(t, mockT,
			"github.com/mattmcnew/archtest/examples/testpackage/nested/dep",
			"github.com/mattmcnew/archtest/examples/testfiledeps/testonlydependency",
		)
	})

	t.Run("Supports checking imports from test packages", func(t *testing.T) {
		mockT := new(testingT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage/...").
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/testfiledeps/testpkgdependency")

		assertNoError(t, mockT)

		archtest.Package(mockT, "github.com/mattmcnew/archtest/examples/testpackage/...").
			IncludeTests().
			ShouldNotDependOn("github.com/mattmcnew/archtest/examples/testfiledeps/testpkgdependency")

		assertError(t, mockT,
			"github.com/mattmcnew/archtest/examples/testpackage/nested/dep_test",
			"github.com/mattmcnew/archtest/examples/testfiledeps/testpkgdependency",
		)
	})
}

func assertNoError(t *testing.T, mockT *testingT) {
	t.Helper()
	if mockT.errored() {
		t.Fatal("archtest should not have failed")
	}
}

func assertError(t *testing.T, mockT *testingT, dependencyTrace ...string) {
	t.Helper()
	if !mockT.errored() {
		t.Fatal("archtest did not fail on dependency")
	}

	if dependencyTrace == nil {
		return
	}

	s := strings.Builder{}
	s.WriteString("Error:\n")
	for i, v := range dependencyTrace {
		s.WriteString(strings.Repeat("\t", i))
		s.WriteString(v + "\n")
	}

	if mockT.message() != s.String() {
		t.Errorf("expected %s got error message: %s", s.String(), mockT.message())
	}
}

type testingT struct {
	errors [][]interface{}
}

func (t *testingT) Error(args ...interface{}) {
	t.errors = append(t.errors, args)
}

func (t testingT) errored() bool {
	if len(t.errors) != 0 {
		return true
	}

	return false

}

func (t *testingT) message() interface{} {
	return t.errors[0][0]
}
