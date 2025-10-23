package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type executorExpected struct {
	returnCode int
	stdOutput  string
}

type executorTestCase struct {
	command  []string
	env      Environment
	input    []string
	expected executorExpected
}

type executorTestSuite struct {
	suite.Suite
	stdOut   *os.File
	stdIn    *os.File
	cmdOut   *os.File
	cmdIn    *os.File
	testCase executorTestCase
}

func (suite *executorTestSuite) SetupTest() {
	suite.cmdOut, _ = os.Create("./testdata/test_output.txt")
	suite.cmdIn, _ = os.Create("./testdata/test_input.txt")
	suite.stdOut = os.Stdout
	suite.stdIn = os.Stdin
	os.Stdout = suite.cmdOut
	os.Stdin = suite.cmdIn
}

func (suite *executorTestSuite) TearDownTest() {
	os.Remove("./testdata/test_output.txt")
	os.Remove("./testdata/test_input.txt")
	os.Stdout = suite.stdOut
	os.Stdin = suite.stdIn
}

func (suite *executorTestSuite) TestSimpleLSCommandWithCheckStdOutput() {
	suite.testCase = executorTestCase{
		command:  []string{"ls", "/bin/bash"},
		expected: executorExpected{returnCode: 0, stdOutput: "/bin/bash\n"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestWithParametersAndCheckStdOutput() {
	suite.testCase = executorTestCase{
		command:  []string{"head", "-n 1", "./testdata/env/BAR"},
		expected: executorExpected{returnCode: 0, stdOutput: "bar\n"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestUnknownCommand() {
	suite.testCase = executorTestCase{
		command:  []string{"this_command_is_unknown"},
		expected: executorExpected{returnCode: 1},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestEnvVariable() {
	suite.testCase = executorTestCase{
		command: []string{"/bin/sh", "-c", "echo $BAR"},
		env: Environment{
			"BAR": EnvValue{
				Value:      "foo",
				NeedRemove: false,
			},
		},
		expected: executorExpected{returnCode: 0, stdOutput: "foo\n"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestRemoveEnvVariable() {
	suite.testCase = executorTestCase{
		command: []string{"/bin/sh", "-c", "echo $HOME"},
		env: Environment{
			"HOME": EnvValue{
				NeedRemove: true,
			},
		},
		expected: executorExpected{returnCode: 0, stdOutput: "\n"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) TestInput() {
	suite.testCase = executorTestCase{
		command:  []string{"head", "-n 1"},
		input:    []string{"This is just a string"},
		expected: executorExpected{returnCode: 0, stdOutput: "This is just a string"},
	}

	suite.runTest()
}

func (suite *executorTestSuite) runTest() {
	suite.writeInput()

	code := RunCmd(suite.testCase.command, suite.testCase.env)

	suite.Equal(suite.testCase.expected.returnCode, code)
	suite.Equal(suite.testCase.expected.stdOutput, suite.readOutput())
}

func (suite *executorTestSuite) readOutput() string {
	suite.cmdOut.Seek(0, 0)
	buf := make([]byte, 1024)
	read, _ := suite.cmdOut.Read(buf)

	return string(buf[0:read])
}

func (suite *executorTestSuite) writeInput() {
	for _, line := range suite.testCase.input {
		suite.cmdIn.Write([]byte(line))
	}
	suite.cmdIn.Seek(0, 0)
}

func TestRunCmd(t *testing.T) {
	suite.Run(t, new(executorTestSuite))
}
