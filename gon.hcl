# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
source = ["./dist/mockthis-macos_darwin_amd64/mockthis", "./dist/mockthis-macos_darwin_arm64/mockthis"]
bundle_id = "com.nicobistolfi.mockthis.cli"

apple_id {
  username = "nbistolfi@gmail.com"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Nico Bistolfi"
}