#include <stdio.h>
#include "allheaders.h"
#include "capi.h"

void die(const char *errstr)
{
	fputs(errstr, stderr);
	exit(1);
}

int main(int argc, char *argv[])
{
	TessBaseAPI *handle;
	PIX *img;
	char *text;

	if ((img = pixRead("img.png")) == NULL)
		die("Error reading image\n");

	handle = TessBaseAPICreate();
	if (TessBaseAPIInit3(handle, NULL, "eng") != 0)
		die("Error initializing tesseract\n");

	TessBaseAPISetImage2(handle, img);
	if (TessBaseAPIRecognize(handle, NULL) != 0)
		die("Error in Tesseract recognition\n");

	if ((text = TessBaseAPIGetUTF8Text(handle)) == NULL)
		die("Error getting text\n");

	TessBaseAPIGetLoadedLanguagesAsVector(handle);

	fputs(text, stdout);

	TessDeleteText(text);
	TessBaseAPIEnd(handle);
	TessBaseAPIDelete(handle);
	pixDestroy(&img);

	return 0;
}