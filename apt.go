package main

import (
  "os/exec"
  "strings"
)

var AptPackageCheck = false

// check for a package on apt

type APTPackage struct {
  Name string
}

func AptSearch(PackageName string) []string {
  if !AptPackageCheck {
    return []string{}
  }
  cmd := exec.Command("apt", "search", PackageName)
  output, err := cmd.Output()
  if err != nil {
    println("Couldn't check the packages")
  }
  pkgs := []string{}
  for _, pkg := range strings.Split(string(output), "\n") {
    if strings.Contains(pkg, "/") && !strings.HasPrefix(pkg, "  "){
      pkgs = append(pkgs, strings.Split(pkg, "/")[0])
    }
  }
  return pkgs
}

func HasAptPackage(PackageName string) bool {
  pkgs := AptSearch(PackageName)
  for _, pkgx := range pkgs {
    if pkgx == PackageName {
      return true
    }
  }
  return false
}