import type { PossibleAuthors } from "@/consts";
import { possibleAuthors } from "@/consts";

export function parseAuthorName(s: string, file: string): PossibleAuthors {

  if (typeof s !== "string") { throw Error("Author name at " + file + " is  not string!") }

  for (const [key, names] of Object.entries(possibleAuthors)) {
    if (names.includes(s.toLowerCase())) {
      console.log(key)
      return key as PossibleAuthors;
    }
  }

  throw Error("Author of name: " + s + " not found")

}

export function parseDateString(s: string, file: string): string {
  const d = new Date(s)

  if (!d) { throw Error("Date from file " + file + " not valid") }

  return d.toLocaleDateString("eu", { month: '2-digit', day: '2-digit', year: 'numeric', })

}
