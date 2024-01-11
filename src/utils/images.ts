import fs from "fs"
import path from "path"

export function getImageDetails(img: string) {
  const imgPath = path.join(process.cwd(), "public", img) 

  const imgSize = fs.statSync(imgPath).size
  const extension = path.extname(imgPath)

  return {
    path: imgPath,
    size: imgSize,
    ext: extension
  }

}
