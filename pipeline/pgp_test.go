package pipeline

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const tempPrefix = "byte-piper-test"

const (
	privatKey = `-----BEGIN PGP PRIVATE KEY BLOCK-----
Version: GnuPG v1

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
=k5D4
-----END PGP PRIVATE KEY BLOCK-----
`
	pubKey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
Version: GnuPG v1

mQENBFRzNl0BCADH/cTVwyaRoxm5uHPAbQ8ltadLkvWgSbOMXM6Z8yeF3njiNxNo
vW0RGFv6XtoRrhUyzaXU92GQ8JORzhRLaKdKcWGlP3Fgi1YD4anhbzpsQ6G2kDeg
P0+GTqGNh6Knz7JnAvrELVjdJGm7YyR848WNtCGsGybJFEnbWyNaThlhYjhhfHec
f2UixIKA5sM2sqfMTUhLjk6XGw0tGK0RlXOUsD2Xw6hFElde6WE/q4Y1Uzx23Uu9
9NBJ1aPqYx+sXXLRGOvjwu3zgCvuPsc/5zw0RIqPWwow98ImlmAxJoUwjYvyW9th
vSk8SxCnhK5HJOGOKuqQYudBRHwGg5hhAMmPABEBAAG0JVRlc3QgJ05vIFRydXN0
JyBNZSA8dGVzdEBleGFtcGxlLmNvbT6JATgEEwECACIFAlRzNl0CGwMGCwkIBwMC
BhUIAgkKCwQWAgMBAh4BAheAAAoJEJF8oPE99cmuesAIAMC9AwHmyj3B2wLZKEib
LyhNuF3U5n9VOrN1UsenjJ+MmLe767jT/hVnqLuCg44fdfZE6ZIEKJdg/LPI3OeO
I0C7CO0+Z5cc8m+grok3v2xps7hcr5XuwIwV11zCMS72Z9R4jB2sztQbToILpCTt
lkKB0YzFv48csbQY+ZRmOBJaI7cIcXleXS4OI1Y7QXmhvYg8MTOD9j2dZ3sLLJDa
WHYJqYRioBEJOT++osT0PZFuMFIj1DiRSvsfOU/lWEArVyE2NYAGaNAGgt24nGxj
Lpzpi7qyntW1vwIJPlxfWy++trMG9ohfMlVA19FSYdgm8hHE8vAMtasOkV/L8t/6
q2u5AQ0EVHM2XQEIALFSuSiP6K1cmf8IFScHAd08qhbGzOeodQQhuG7WRWHVlNm2
RB3VLrz15NqtinkjCiTJ/8JbAoDQzKvB72pjlpy0Rwn/ZFndMcK2YG+6TrgfToSU
DXqSXJRGlaXZd6yBrWt4mU6v7QTwaGH7LhLNTniCoatS6+NpUwynbTHfrTR2Knmu
V1z8j06cpkm50STdmKF9EfPv15+NlQGxIMNq+dg5rrk47kpyWVugvGz95epPEkZn
w8k8cVnFULoEh5PaAwPwVjC1dCHt2rbP2SNiwZen8wdl7r8Y/XMpV3BFr2ehOK0l
T82V5w1pvJB/C1B6bIkYpA8gOTcSRE7kfCI5ZO0AEQEAAYkBHwQYAQIACQUCVHM2
XQIbDAAKCRCRfKDxPfXJrnyVB/9biEfpIUPCg5WegDdkxZEeMjEIg4Km/tS5aCmK
SpO2D2fo6X3LfyicVMedFD4AUD3dNSpbnnnkpiQTUrQE++TFX9bBv39ciYIMpEYo
82u35SkVpnLVm8Kgvi1euqlzf/EVkHXJp6m+h3G4xRtBGDvILpE2+oebu3q7ev1J
eQBVKOnLTofNmwiDL3+H3zx3m8QpX9nb2MynodhuWvAMp9jaCi0VaWRh1zijqZ1M
xdYI2Eqo+Hzv+jHCal0nqc4o7rEKnNflrAhraN3+//FJFiOp6nYeItOsOHUeQCmq
zo4z737LeL7ekrNgI8g/6+93KKMNL6j2nzfSvMCtpvpppeaY
=i8mq
-----END PGP PUBLIC KEY BLOCK-----
`
	keyId        = "3DF5C9AE"
	expectedText = `ἔστι δὲ ἡ σκυτάλη τοιοῦτον. ἐπὰν ἐκπέµπωσι ναύαρχον ἢ
στρατηγὸν οἱ ἔφοροι, ξύλα δύο στρογγύλα µῆκος καὶ πάχος ἀκριβῶς
ἀπισώσαντες, ὥστε ταῖς τοµαῖς ἐφαρµόζειν πρὸς ἄλληλα, τὸ µὲν αὐτοὶ
φυλάττουσι, θάτερον δὲ τῷ πεµποµένῳ διδόασι. ταῦτα δὲ τὰ ξύλα
σκυτάλας καλοῦσιν. ὅταν οὖν ἀπόρρητόν τι καὶ µέγα φράσαι
βουληθῶσι, βιβλίον ὥσπερ ἱµάντα µακρὸν καὶ στενὸν ποιοῦντες
περιελίττουσι τὴν παρ' αὐτοῖς σκυτάλην, οὐδὲν διάλειµµα ποιοῦντες, 
ἀλλὰ πανταχόθεν κύκλῳ τὴν ἐπιφάνειαν αὐτῆς τῷ βιβλίῳ
καταλαµβάνοντες. τοῦτο δὲ ποιήσαντες ἃ βούλονται καταγράφουσιν εἰς
τὸ βιβλίον, ὥσπερ ἐστὶ τῇ σκυτάλῃ περικείµενον· ὅταν δὲ γράψωσιν, 
ἀφελόντες τὸ βιβλίον ἄνευ τοῦ ξύλου πρὸς τὸν στρατηγὸν ἀποστέλλουσι. 
δεξάµενος δὲ ἐκεῖνος ἄλλως µὲν οὐδὲν ἀναλέξασθαι δύναται τῶν
γραµµάτων συναφὴν οὐκ ἐχόντων, ἀλλὰ διεσπασµένων, τὴν δὲ παρ' αὑτῷ
σκυτάλην λαβὼν τὸ τµῆµα τοῦ βιβλίου περὶ αὐτὴν περιέτεινεν, ὥστε, τῆς
ἕλικος εἰς τάξιν ὁµοίως ἀποκαθισταµένης, ἐπιβάλλοντα τοῖς πρώτοις τὰ
δεύτερα, κύκλῳ τὴν ὄψιν ἐπάγειν τὸ συνεχὲς ἀνευρίσκουσαν. καλεῖται δὲ
ὁµωνύµως τῷ ξύλῳ σκυτάλη τὸ βιβλίον, ὡς τῷ µετροῦντι τὸ µετρούµενον.`
)

func TestPGP(t *testing.T) {
	in := bytes.NewBuffer([]byte(expectedText))
	out := &bytes.Buffer{}

	pgp, err := newPGPFilter(map[string]string{"pubkey": pubKey})
	if err != nil {
		t.Fatal(err)
	}
	if err := pgp.Link(in); err != nil {
		t.Fatal(err)
	}

	unpgp, err := newUnpgpFilter(map[string]string{"privatkey": privatKey})
	if err != nil {
		t.Fatal(err)
	}
	if err := unpgp.Link(pgp); err != nil {
		t.Fatal(err)
	}

	if _, err := io.Copy(out, unpgp); err != nil {
		t.Fatal(err)
	}

	if out.String() != expectedText {
		t.Fatal("Unexpected string %s", out.String())
	}
}

func TestPGPforGPG(t *testing.T) {
	in := bytes.NewBuffer([]byte(expectedText))

	pgp, err := newPGPFilter(map[string]string{"pubkey": pubKey})
	if err != nil {
		t.Fatal(err)
	}
	if err := pgp.Link(in); err != nil {
		t.Fatal(err)
	}

	gpgHome, err := ioutil.TempDir("", tempPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(gpgHome)

	if err := gpgImportKey(gpgHome, privatKey); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("gpg", "--homedir", gpgHome, "-d")
	cmd.Stdin = pgp
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	cmd.Wait()
	if string(out) != expectedText {
		t.Fatalf("Unexpected string %s", out)
	}

}

func TestGPGforPGP(t *testing.T) {
	gpgHome, err := ioutil.TempDir("", tempPrefix)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(gpgHome)

	unpgp, err := newUnpgpFilter(map[string]string{"privatkey": privatKey})
	if err != nil {
		t.Fatal(err)
	}

	if err := gpgImportKey(gpgHome, pubKey); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("gpg", "--homedir", gpgHome, "-er", keyId, "--trust-model", "always")
	cmd.Stdin = strings.NewReader(expectedText)
	stdout, err := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := unpgp.Link(stdout); err != nil {
		t.Fatal(err)
	}
	out, err := ioutil.ReadAll(unpgp)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != expectedText {
		t.Fatalf("Unexpected string '%s'", out)
	}
}

func gpgImportKey(home, key string) error {
	cmd := exec.Command("gpg", "--homedir", home, "--import")
	cmd.Stdin = strings.NewReader(key)
	return cmd.Run()
}
