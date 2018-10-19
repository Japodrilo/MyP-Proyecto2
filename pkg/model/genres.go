package model

// Genre is a simple translator for ID3v1 genre codes, it contains
// a single dictionary with the codes as keys and the corresponding
// genres as values (both are strings).  It is a singleton.
type Genre struct {
    genres map[string]string
}

var instance *Genre

// GetGenre returns the singleton instance of Genre. The implementation
// of the singleton is not thread-safe due to conflicts with gtk.
func GetGenre() *Genre {
    if instance == nil {
        genres := make(map[string]string)
        genres["0"] = "Blues"
        genres["1"] = "Classic Rock"
        genres["2"] = "Country"
        genres["3"] = "Dance"
        genres["4"] = "Disco"
        genres["5"] = "Funk"
        genres["6"] = "Grunge"
        genres["7"] = "Hip-Hop"
        genres["8"] = "Jazz"
        genres["9"] = "Metal"
        genres["10"] = "New Age"
        genres["11"] = "Oldies"
        genres["12"] = "Other"
        genres["13"] = "Pop"
        genres["14"] = "R&B"
        genres["15"] = "Rap"
        genres["16"] = "Reggae"
        genres["17"] = "Rock"
        genres["18"] = "Techno"
        genres["19"] = "Industrial"
        genres["20"] = "Alternative"
        genres["21"] = "Ska"
        genres["22"] = "Death Metal"
        genres["23"] = "Pranks"
        genres["24"] = "Soundtrack"
        genres["25"] = "Euro-Techno"
        genres["26"] = "Ambient"
        genres["27"] = "Trip-Hop"
        genres["28"] = "Vocal"
        genres["29"] = "Jazz+Funk"
        genres["30"] = "Fusion"
        genres["31"] = "Trance"
        genres["32"] = "Classical"
        genres["33"] = "Instrumental"
        genres["34"] = "Acid"
        genres["35"] = "House"
        genres["36"] = "Game"
        genres["37"] = "Sound Clip"
        genres["38"] = "Gospel"
        genres["39"] = "Noise"
        genres["40"] = "AlternRock"
        genres["41"] = "Bass"
        genres["42"] = "Soul"
        genres["43"] = "Punk"
        genres["44"] = "Space"
        genres["45"] = "Meditative"
        genres["46"] = "Instrumental Pop"
        genres["47"] = "Instrumental Rock"
        genres["48"] = "Ethnic"
        genres["49"] = "Gothic"
        genres["50"] = "Darkwave"
        genres["51"] = "Techno-Industrial"
        genres["52"] = "Electronic"
        genres["53"] = "Pop-Folk"
        genres["54"] = "Eurodance"
        genres["55"] = "Dream"
        genres["56"] = "Southern Rock"
        genres["57"] = "Comedy"
        genres["58"] = "Cult"
        genres["59"] = "Gangsta"
        genres["60"] = "Top 40"
        genres["61"] = "Christian Rap"
        genres["62"] = "Pop/Funk"
        genres["63"] = "Jungle"
        genres["64"] = "Native American"
        genres["65"] = "Cabaret"
        genres["66"] = "New Wave"
        genres["67"] = "Psychedelic"
        genres["68"] = "Rave"
        genres["69"] = "Showtunes"
        genres["70"] = "Trailer"
        genres["71"] = "Lo-Fi"
        genres["72"] = "Tribal"
        genres["73"] = "Acid Punk"
        genres["74"] = "Acid Jazz"
        genres["75"] = "Polka"
        genres["76"] = "Retro"
        genres["77"] = "Musical"
        genres["78"] = "Rock & Roll"
        genres["79"] = "Hard Rock"
        instance = &Genre{genres}
    }
    return instance
}

// Get receives a string corresponding to an ID3v1 genre code
// and returns the corresponding genre.
func (genre *Genre) Get(code string) string {
    name := genre.genres[code]
    if name == "" {
        return code
    }
    return name
}
