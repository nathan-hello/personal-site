package utils

import "strings"

const AsciiCommentDesc = "Ok so there are a couple things you can do: \n    - Make an arse of yourself (don't do this) \n    - Use markdown styling like codeblocks by delimiting \n      with backticks: \n        ```C \n        int main(void) { \n            printf(\"Hello, World\\n\"); \n            return 0; \n        } \n        ``` \n    - Write inline latex with dollars: \n        $R_{\\mu \\nu} - \\frac{1}{2} R \\, g_{\\mu \\nu} \n        + \\Lambda g_{\\mu \\nu} = \\frac{8 \\pi G}{c^4} \\, T_{\\mu \\nu}$ \n    - basic green text and post linking with > and >> respectively \n    All rendered without any client side js."

const AsciiNat_e = `
 _   _       _       __   
| \ | |     | |     / /   
|  \| | __ _| |_   / /__  
| . ! |/ _! | __| / / _ \ 
| |\  | (_| | |_ / /  __/ 
|_| \_|\__,_|\__/_/ \___| "We are Boingus"
`

const AsciiIrc = `
     __                    __
    / / ________  ______  / /
   / / /  _/ __ \/ ____/ / /
  / /  / // /_/ / /     / / 
 / / _/ // _, _/ /___  / /
/_/ /___/_/ |_|\____/ /_/ 
`

const AsciiNatzone = `
                                                            _______         .--.    
      __                                         __       _ \______ -      |o_o |
     / /_   __      __                          / /      | \  ___  \ |     |:_/ |
    / // | / /___ _/ /_____  ____  ____  ___   / /       | | /   \ | |    //   \ \
   / //  |/ / __ '/ __/_  / / __ \/ __ \/ _ \ / /        | | \___/ | |   (|     | )
  / // /|  / /_/ / /_  / /_/ /_/ / / / /  __// /         | \______ \_|  /'|_   _/'\
 /_//_/ |_/\__,_/\__/ /___/\____/_/ /_/\___//_/           -_______\     \___)=(___/
`

const AsciiLinks = `
      __                        __
     / / ___       __          / /
    / / / (_)___  / /_______  / / 
   / / / / / __ \/ //_/ ___/ / /  
  / / / / / / / / , (__  )  / /   
 /_/ /_/_/_/ /_/_/|_/____/ /_/      
`

const AsciiTv = `
      __              __
     / / __          / /
    / / / /__   __  / / 
   / / / __/ | / / / /  
  / / / /_ | |/ / / /   
 /_/  \__/ |___/ /_/
`

const AsciiWebring = `
      __                __         _                __   
     / / _      _____  / /_  _____(_)___  ____ _   / / 
    / / | | /| / / _ \/ __ \/ ___/ / __ \/ __ '/  / /  
   / /  | |/ |/ /  __/ /_/ / /  / / / / / /_/ /  / /   
  / /   |__/|__/\___/_.___/_/  /_/_/ /_/\__, /  / /    
 /_/                                   /____/  /_/
`

func AsciiRender(s string) string {
        return strings.ReplaceAll(s, "!", "`")
}
