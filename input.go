package main

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
)

type UserInput struct {
	reader UiReader
}

func (u *UserInput) Ask(q string, retries int, v validator, vMsg string) (string, error) {
	var text string
	var err error
	print(q)
	for i := 0; i < retries; i++ {

		text, err = u.getInput()

		if v(text) && err == nil {
			return text, nil
		}
		println(vMsg)

		print(q)
	}

	return "", errors.New("problem getting user input")
}

func (u *UserInput) AskYesOrNo(q string, retries int) (bool, error) {
	text, err := u.Ask(q, retries, isYorN, "(please enter a single letter 'y' or 'n')")

	if err != nil {
		return false, err
	}

	if strings.ToLower(text[0:1]) == "y" {
		return true, nil
	}
	return false, nil
}

func (u *UserInput) AskWithOptions(q string, opts []string, retries int) (int, error) {
	if len(opts) == 0 {
		return 0, errors.New("no options")
	}
	var s int
	println(q)

	for i, opt := range opts {
		println(" " + strconv.Itoa(i+1) + ". " + opt)
	}

	print("Enter your selection: ")
	for i := 0; i < retries; i++ {

		text, err := u.getInput()

		if err == nil {
			s, err = strconv.Atoi(text)
			if err == nil {
				return s, nil
			}
		}

		print("Invalid selection (" + err.Error() + " ). Please enter a number between 1 and " + strconv.Itoa(len(opts)) + ": ")
	}

	return 0, errors.New("user failed input")
}

type validator func(input string) (valid bool)

type UiReader interface {
	ReadString(byte) (string, error)
}

func (u *UserInput) getInput() (string, error) {
	text, err := u.reader.ReadString('\n')
	if err != nil {
		return text, err
	}
	return strings.Join(strings.Fields(text), ""), nil
}

type BufIoWrapper struct {
	b bufio.Reader
}

func (w BufIoWrapper) ReadString(term byte) (string, error) {
	return w.b.ReadString(term)
}

func isYorN(text string) bool {
	l := strings.ToLower(text)
	if l == "y" || l == "n" || l == "yes" || l == "no" {
		return true
	}
	return false
}
