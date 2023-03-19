// #EXTM3U - It is the file header indicating Extended M3U and must be first line of the file.
// #EXTENC: - Text encoding. It must be the 2nd line of the file.
// #EXTINF: - Used for track information and other additional properties.
// #PLAYLIST: - The title of the playlist
// #EXTGRP: - Begin name grouping
// #EXTALB: - Album information
// #EXTART: - Album artist
// #EXTGENRE - Album Genre
// #EXTM3A - Single file playlist for album tracks or chapters.
// #EXTBYT: - File size in bytes.
// #EXTBIN: - Binary data follows.
// #EXTIMG: - Logo, Cover or other images.

// #EXTM3U8

// #EXTINF:111, Sample artist name - Sample track title
// C:\Music\SampleMusic.mp3

// https://www.rfc-editor.org/rfc/rfc8216

package utils

import "fmt"

func ToM3u8(metas []map[string]interface{}, url string) string {
	result := "#EXTM3U\n\n"

	for _, meta := range metas {
		result += "#EXTINF:" + fmt.Sprintf("%v, %v", meta["duration"], meta["title"]) + "\n" +
			"#EXTALB:" + fmt.Sprintf("%v", meta["album"]) + "\n" +
			"#EXTART:" + fmt.Sprintf("%v", meta["artist"]) + "\n" +
			url + "\n\n"
	}

	return result + "#EXT-X-ENDLIST\n"
}
