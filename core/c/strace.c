#include <stdio.h>
#include <stdlib.h>

int main(void) {
	FILE *f;
	char str[1024];

	if (!(f=fopen("/etc/passwd", "r"))) {
		fprintf(stderr, "Could not open file");
		exit(1);
	}

	while(!feof(f)) {
		fscanf(f, "%s", str);
		fprintf(stdout, "%s", str);
	}

	exit(0);
}
