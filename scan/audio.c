#include "stdafx.h"

static int64_t toLow(char *str)
{
    int64_t _buf = 0;
    char *buf = (char *)(&_buf);

    for (int8_t i = 0; i < 7 && str[i]; i++)
    {
        char ch = str[i];
        buf[i] = ('A' <= ch && ch <= 'Z') ? (ch + 32) : ch;
    }

    return _buf;
}

static AVCodecContext *new_codec_cxt(const AVCodec *codec)
{
    AVCodecContext *coder_ctx;

    if (!(coder_ctx = avcodec_alloc_context3(codec)))
        return NULL;

    if (avcodec_open2(coder_ctx, codec, NULL) < 0)
    {
        avcodec_free_context(&coder_ctx);
        return NULL;
    }

    return coder_ctx;
}

/**`
 * @return {
 *   cover: string; // https://img.webp
 *   lrc: string; // lrc txt
 *   url: string; // /sha1.flac
 *
 *   title: string; // My all
 *   artist: string; // 浜崎あゆみ
 *   album: string; // GUILTY
 *
 *   smple_rate: number; // 44100(Hz)
 *   bit_rate: number; // 1077(kbps)
 *   channles: number; // 2
 * }
 */
int load_audio(const char *inputFileName, void *buf, size_t size)
{
    AVFormatContext *fmt_ctx = NULL;
    AVDictionaryEntry *tag = NULL;
    const AVCodec *codec = NULL;
    AVCodecContext *coder_ctx = NULL;
    int ret = -1;

    if (avformat_open_input(&fmt_ctx, inputFileName, NULL, NULL))
        return -1;

    if (avformat_find_stream_info(fmt_ctx, NULL) < 0)
        goto End;

    const int stream_nb = av_find_best_stream(fmt_ctx, AVMEDIA_TYPE_AUDIO, -1, -1, &codec, 0);
    if (stream_nb < 0)
        goto End;

    AVStream *av_stream = fmt_ctx->streams[stream_nb];

    if (!(coder_ctx = new_codec_cxt(codec)))
    {
        // ret = AVERROR(ENOMEM);
        goto End;
    }

    if (avcodec_parameters_to_context(coder_ctx, av_stream->codecpar) < 0)
        goto End;

    size_t len = 0;
    // printf("path=%s\n", inputFileName);
    len += snprintf(buf, size, "{\"sample_rate\":%d,", coder_ctx->sample_rate);             // Hz
    len += snprintf(buf + len, size - len, "\"bit_rate\":%ld,", fmt_ctx->bit_rate / 1000); // kbps
    len += snprintf(buf + len, size - len, "\"channels\":%d,", coder_ctx->ch_layout.nb_channels);
    while ((tag = av_dict_get(fmt_ctx->metadata, "", tag, AV_DICT_IGNORE_SUFFIX)))
    {
        int64_t key = toLow(tag->key);
        // 0x6D75626C61, album
        // 0x747369747261, artist
        // 0x656C746974, title
        if (0x6D75626C61 == key || 0x747369747261 == key || 0x656C746974 == key)
            len += snprintf(buf + len, size - len, "\"%s\": \"%s\",", (char *)(&key), tag->value);
    }
    len += snprintf(buf + len, size - len, "\"duration\":%ld}", fmt_ctx->duration / AV_TIME_BASE);

    ret = len;

End:
    if (NULL != coder_ctx)
        avcodec_free_context(&coder_ctx);
    // printf("%ld\n", &fmt_ctx);
    avformat_close_input(&fmt_ctx);
    // printf("%d\n", 321321);
    return ret;
}
