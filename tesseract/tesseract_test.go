package tesseract

import (
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

// func TestGetText(t *testing.T) {
// 	handle, err := TessBaseAPICreate("tur")
// 	if err != nil {
// 		t.Fatalf("err creating tesseract: %v", err)
// 	}

// 	testCases := []struct {
// 		desc string
// 		path string
// 	}{
// 		{
// 			desc: "a.png",
// 			path: "../testdata/a.png",
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := handle.SetImage2(tC.path)
// 			if err != nil {
// 				t.Errorf(err.Error())
// 			}
// 			err = handle.Recognize()
// 			if err != nil {
// 				t.Errorf(err.Error())
// 			}

// 			text, err := handle.GetUTF8Text()
// 			if err != nil {
// 				t.Errorf(err.Error())
// 			}
// 			t.Logf("got the text: %s", text)
// 		})
// 	}
// }

// func TestGetTextFromTIFF(t *testing.T) {
// 	handle, err := TessBaseAPICreate("tur")
// 	if err != nil {
// 		t.Fatalf("err creating tesseract: %v", err)
// 	}

// 	testCases := []struct {
// 		desc string
// 		path string
// 	}{
// 		{
// 			desc: "p2.tif",
// 			path: "../trashdata/p2.tif",
// 		},
// 	}
// 	for _, tC := range testCases {
// 		t.Run(tC.desc, func(t *testing.T) {
// 			err := handle.SetImage2(tC.path)
// 			if err != nil {
// 				t.Errorf(err.Error())
// 			}
// 			err = handle.Recognize()
// 			if err != nil {
// 				t.Errorf(err.Error())
// 			}

// 			text, err := handle.GetUTF8Text()
// 			if err != nil {
// 				t.Errorf(err.Error())
// 			}
// 			t.Logf("got the text: %s", text)
// 		})
// 	}
// }

func TestUsePixSlice(t *testing.T) {
	handle, err := TessBaseAPICreate("tur")
	if err != nil {
		t.Fatalf("err creating tesseract: %v", err)
	}

	testCases := []struct {
		desc string
		path string
	}{
		{
			desc: "p2.tif",
			path: "./testdata/multipage-sample.tif",
		},
		{
			desc: "a.png",
			path: "./testdata/a.png",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			text, err := handle.ProcessImage(tC.path)
			if err != nil {
				t.Errorf(err.Error())
			}
			t.Logf("got the text: %v", text)
		})
	}
}
