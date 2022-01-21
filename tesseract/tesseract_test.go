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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
