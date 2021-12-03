#include "sum_array.h"

int32_t sum_array(uint32_t size, int32_t *array)
{
	int32_t sum = 0;
	for (uint32_t i = 0; i < size; i++)
	{
		sum += array[i];
	}

	return sum;
}
