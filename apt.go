package main

import "os/exec"

// check for a package on apt

type APTPackage struct {
  Maintainer string
  Name string
  Description string
}

func APTSearch(PackageName string) {
  cmd := exec.Command("apt", "search", PackageName)
  output, err := cmd.Output()
  if err != nil {
    println("Couldn't check the packages")
  }
  println(string(output))
}