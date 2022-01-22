package tesseract

import (
	"reflect"
	"testing"
)

func TestTessVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1st",
			want: "4.1.3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TessVersion(); got != tt.want {
				t.Errorf("TessVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
	images in testdata were taken from the repo below
	https://github.com/madmaze/pytesseract/tree/master/tests/data

	i wrote the repo below to be able to download the files
	inside of a directory from github directly
	if you are interested in
	https://github.com/rdmyldz/gitapi
*/
func TestTessBaseAPI_ProcessImage(t *testing.T) {
	const testMessage string = "This is test message\n"
	const testLazyFox string = `This is a lot of 12 point text to test the
ocr code and see if it works on all types
of file format.

The quick brown dog jumped over the
lazy fox. The quick brown dog jumped
over the lazy fox. The quick brown dog
jumped over the lazy fox. The quick
brown dog jumped over the lazy fox.
`
	const testEuropean string = `The (quick) [brown] {fox} jumps!
Over the $43,456.78 <lazy> #90 dog
& duck/goose, as 12.5% of E-mail
from aspammer@website.com is spam.
Der ,.schnelle” braune Fuchs springt
iiber den faulen Hund. Le renard brun
«rapide» saute par-dessus le chien
paresseux. La volpe marrone rapida
salta sopra il cane pigro. El zorro
marron rapido salta sobre el perro
perezoso. A raposa marrom rapida
salta sobre o cdo preguicoso.
`

	// note: don't have .pnm and .spx files, jp2 is already not supported
	tests := []struct {
		name    string
		path    string
		lang    string
		want    []string
		wantErr bool
	}{
		{
			name:    "nonexistent image",
			path:    "testdata/nonexist.png",
			lang:    "eng",
			wantErr: true,
		},
		{
			name:    "test_la.png",
			path:    "testdata/test_la.png",
			lang:    "eng",
			want:    []string{testMessage},
			wantErr: false,
		},
		{
			name:    "test.png",
			path:    "testdata/test.png",
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test.tiff",
			path:    "testdata/test.tiff",
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test.bmp",
			path:    "testdata/test.bmp",
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test.jpg",
			path:    "testdata/test.jpg",
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test.gif",
			path:    "testdata/test.gif",
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test.webp",
			path:    "testdata/test.webp",
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test.pgm",
			path:    "testdata/test.pgm", // leptonica.IFF_PNM
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test.ppm",
			path:    "testdata/test.ppm", // leptonica.IFF_PNM
			lang:    "eng",
			want:    []string{testLazyFox},
			wantErr: false,
		},
		{
			name:    "test-small.jpg",
			path:    "testdata/test-small.jpg",
			lang:    "eng",
			want:    []string{"This\n"},
			wantErr: false,
		},
		{
			name:    "test-european.jpg",
			path:    "testdata/test-european.jpg",
			lang:    "eng",
			want:    []string{testEuropean},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handle, err := TessBaseAPICreate(tt.lang)
			if err != nil {
				t.Fatalf("error creating handle: %v", err)
			}
			got, err := handle.ProcessImage(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("handle.ProcessImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handle.ProcessImage() = %q, want %v", got, tt.want)
			}
		})
	}
}
