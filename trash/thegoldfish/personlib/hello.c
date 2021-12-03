#include <stdio.h>
#include "person.h"

int main(int argc, char const *argv[])
{
	APerson *of;
	of = get_person("tim", "tim hughes");
	printf("Hello C world: My name is %s, %s.\n", of->name, of->long_name);

	return 0;
}
