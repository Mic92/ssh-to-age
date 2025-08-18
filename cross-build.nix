{ lib, pkgs, ssh-to-age }:

let
  # Define all target platforms for cross-compilation
  platforms = [
    { os = "darwin"; arch = "amd64"; goarch = "amd64"; }
    { os = "darwin"; arch = "arm64"; goarch = "arm64"; }
    { os = "freebsd"; arch = "386"; goarch = "386"; }
    { os = "freebsd"; arch = "amd64"; goarch = "amd64"; }
    { os = "freebsd"; arch = "arm"; goarch = "arm"; }
    { os = "linux"; arch = "386"; goarch = "386"; }
    { os = "linux"; arch = "amd64"; goarch = "amd64"; }
    { os = "linux"; arch = "arm"; goarch = "arm"; }
    { os = "linux"; arch = "arm64"; goarch = "arm64"; }
    { os = "linux"; arch = "mips"; goarch = "mips"; }
    { os = "linux"; arch = "mips64"; goarch = "mips64"; }
    { os = "linux"; arch = "mips64le"; goarch = "mips64le"; }
    { os = "linux"; arch = "mipsle"; goarch = "mipsle"; }
    { os = "linux"; arch = "ppc64"; goarch = "ppc64"; }
    { os = "linux"; arch = "ppc64le"; goarch = "ppc64le"; }
    { os = "linux"; arch = "riscv64"; goarch = "riscv64"; }
    { os = "linux"; arch = "s390x"; goarch = "s390x"; }
    { os = "netbsd"; arch = "386"; goarch = "386"; }
    { os = "netbsd"; arch = "amd64"; goarch = "amd64"; }
    { os = "netbsd"; arch = "arm"; goarch = "arm"; }
    { os = "openbsd"; arch = "386"; goarch = "386"; }
    { os = "openbsd"; arch = "amd64"; goarch = "amd64"; }
    { os = "openbsd"; arch = "arm"; goarch = "arm"; }
    { os = "openbsd"; arch = "arm64"; goarch = "arm64"; }
    { os = "windows"; arch = "386"; goarch = "386"; }
    { os = "windows"; arch = "amd64"; goarch = "amd64"; }
    { os = "windows"; arch = "arm64"; goarch = "arm64"; }
  ];

  buildForPlatform = platform:
    ssh-to-age.overrideAttrs (old: {
      env = (old.env or {}) // {
        CGO_ENABLED = "0";
        GOOS = platform.os;
        GOARCH = platform.goarch;
      };
      doCheck = false;
    });

in
pkgs.runCommand "ssh-to-age-cross-build" { } ''
  mkdir -p $out
  
  # Build for each platform and rename binaries
  ${lib.concatMapStringsSep "\n" (platform:
    let
      suffix = if platform.os == "windows" then ".exe" else "";
      outputName = "ssh-to-age.${platform.os}-${platform.arch}${suffix}";
      build = buildForPlatform platform;
    in
    ''
      echo "Copying ${platform.os}/${platform.arch} binary..."
      # Native binaries are directly in bin/, cross-compiled ones are in subdirectories
      if [ -f "${build}/bin/ssh-to-age${suffix}" ]; then
        cp ${build}/bin/ssh-to-age${suffix} $out/${outputName}
      else
        cp ${build}/bin/${platform.os}_${platform.goarch}/ssh-to-age${suffix} $out/${outputName}
      fi
      sha256sum $out/${outputName} > $out/${outputName}.sha256
    ''
  ) platforms}
  
  # Create combined checksums file with standard name
  cd $out
  cat *.sha256 > sha256sums.txt
  rm -f *.sha256
''
