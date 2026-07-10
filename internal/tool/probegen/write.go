package main

import "os"

// writeGenerated writes a fully generated file. The lone gosec suppression for
// generated, git-committed, world-readable output lives here so no other writer
// repeats it.
func writeGenerated(path, content string) error {
	//nolint:gosec // generated artifact, committed to git, world-readable is fine
	return os.WriteFile(path, []byte(content), 0o644)
}

// replaceMarkedBlock rewrites the text between begin and end markers in the file
// at path. Reading a caller-controlled generator path is safe.
func replaceMarkedBlock(path, begin, end, body string) error {
	//nolint:gosec // caller-controlled generator path, not user input
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	out, err := replaceBlock(string(data), begin, end, body)
	if err != nil {
		return err
	}
	return writeGenerated(path, out)
}
