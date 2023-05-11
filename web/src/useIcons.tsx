import * as icons from "@mui/icons-material"
import stringSimilarity from 'string-similarity'

function useIcons(word) {
    const iconsNames = Object.keys(icons)

    var matches = stringSimilarity.findBestMatch(word, iconsNames)
    const bestMathch = matches.bestMatch.target
    const Icon = icons[bestMathch]
    return Icon
}
export default useIcons