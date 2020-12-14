package builder

import (
	"context"
	"io"

	"github.com/spf13/cobra"
)

type Builder interface {
	WithDescription(description string) Builder
	WithLongDescription(longDescription string) Builder
	SetFlags() Builder
	SetExample() Builder
	ExactArgs(argCount int, action func(context.Context, io.Writer, []string) error) *cobra.Command
	NoArgs(action func(context.Context, io.Writer) error) *cobra.Command
}

type builder struct {
	cmd cobra.Command
}

func NewCmd(use string) Builder {
	return &builder{
		cmd: cobra.Command{
			Use: use,
		},
	}
}

func (b builder) WithLongDescription(longDescription string) Builder {
	b.cmd.Long = longDescription
	return b
}

func (b builder) WithDescription(description string) Builder {
	b.cmd.Short = description
	return b
}

func (b builder) SetExample() Builder {
	SetCommandExample(&b.cmd)
	return b
}

func (b builder) SetFlags() Builder {
	SetCommandFlags(&b.cmd)
	return b
}

func (b builder) WithUsageTemplate(s string) {
	b.cmd.SetUsageTemplate("abc")
}

func (b builder) ExactArgs(argCount int, action func(context.Context, io.Writer, []string) error) *cobra.Command {
	b.cmd.Args = cobra.ExactArgs(argCount)
	b.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return returnErrorFromFunction(action(b.cmd.Context(), b.cmd.OutOrStdout(), args))
	}
	return &b.cmd
}

func (b builder) NoArgs(action func(context.Context, io.Writer) error) *cobra.Command {
	b.cmd.Args = cobra.NoArgs
	b.cmd.RunE = func(*cobra.Command, []string) error {
		return returnErrorFromFunction(action(b.cmd.Context(), b.cmd.OutOrStdout()))
	}
	return &b.cmd
}

func returnErrorFromFunction(err error) error {
	return err
}