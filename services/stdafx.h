#ifndef _STDAFX_H_
#define _STDAFX_H_

#include <stdio.h>
#include <stdlib.h>
#include <libavformat/avformat.h>
#include <libavutil/dict.h>
#include <libavcodec/avcodec.h>

int load_audio(const char *, void *, size_t);

#endif
