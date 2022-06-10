{ lib, buildGoModule, fetchFromGitHub , makeWrapper }:

buildGoModule rec {
  pname = "go-delve";
  version = "1.8.3";

  src = fetchFromGitHub {
    owner = "go-delve";
    repo = "delve";
    rev = "v${version}";
    sha256 = "sha256-6hiUQNUXpLgvYl/MH+AopIzwqvX+vtvp9GDEDmwlqek=";
  };

  vendorSha256 = null;

  CGO_ENABLED = 0;
  # Tests use network
  doCheck = false;
  subPackages = [ "cmd/dlv" ];
  allowGoReference = true;
  checkFlags = [ "-short"];

  nativeBuildInputs = [ makeWrapper ];

  postInstall = ''
    # fortify source breaks build since delve compiles with -O0
    wrapProgram $out/bin/dlv \
      --prefix disableHardening " " fortify
    # add symlink for vscode golang extension
    # https://github.com/golang/vscode-go/blob/master/docs/debugging.md#manually-installing-dlv-dap
    ln $out/bin/dlv $out/bin/dlv-dap
  '';
}
