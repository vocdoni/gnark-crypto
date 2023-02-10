package mimc

import (
	"os"
	"path/filepath"

	"github.com/consensys/bavard"
	"github.com/consensys/gnark-crypto/internal/generator/config"
)

func Generate(conf config.Curve, baseDir string, bgen *bavard.BatchGenerator) error {

	conf.Package = "mimc"
	entries := []bavard.Entry{
		{File: filepath.Join(baseDir, "doc.go"), Templates: []string{"doc.go.tmpl"}},
		{File: filepath.Join(baseDir, "mimc.go"), Templates: []string{"mimc.go.tmpl"}},
		{File: filepath.Join(baseDir, "mimc_test.go"), Templates: []string{"mimc.test.go.tmpl"}},
	}
	os.Remove(filepath.Join(baseDir, "utils.go")) // TODO: Safe to remove these now?
	os.Remove(filepath.Join(baseDir, "utils_test.go"))

	return bgen.Generate(conf, conf.Package, "./crypto/hash/mimc/template", entries...)

}
