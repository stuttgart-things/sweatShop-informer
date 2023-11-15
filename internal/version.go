/*
Copyright © 2023 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"github.com/fatih/color"
	goVersion "go.hein.dev/go-version"
)

var (
	date    = "unknown"
	commit  = "unknown"
	output  = "yaml"
	version = "unset"
)

// https://ascii-generator.site/
const logo = `

                       ...                    .==.
                       .. .                :==-  +=
                                        :++:      :*.
                                  :-==++-    .:.    *=:
                 :%%*+=:     :++**+=##%=.    *%##*=. +%%%++=.
                 #==#*==**=*##%=----+%**##+-  .+%#*####@=--=#%*=..::::..
                .%---+#+--*%**#%=----%**%####*-  +%***%*--=*%#%#*++=+##%#+
                =*----=#+--+###%#****%%#%@+-==+:  ##*#%--=%%#+=-=***+=--=#
                =*-----=%=----------====+++*-   =********##+--=*#+------+*
                +#****++++===---------------+*+*=------------=%=--------#:
               *=     .::--====++++++********+++++===============+++***#%
              :%***+++===--:::..             ..::::-----======--::..    #:
    :....:.   +#+*****+++++++*********++++++=======---------------==++**%=
   -   .:.    #:    .::---====+++************++++++++++++++++++++++****++#
   .:.-       %***+++==--::.            ....:::-:--------======---:.     @
             #+----=%-===+++*****+++++====--::...               .::-===+*+
            #*+++++#*-----+*##%%%####****++***########**********+#*==---+#
           +%######@=-------#*@@@@@@@%  %-=%%@@@@@@@=-%**+-------%=======%:
          -%*******@-------=%#@@@@@@@@. %-**@@@@@@@@= *+--------=@#######%*
         .%********%-------++%@@@@@@@@. %-%+@@@@@@@@= +*--------*%********%.
         %#********%-------*++@@@@@@@# .%-%=@@@@@@@@: *+--------##********#*       .
        +%#########@=------+# #@@@@@*. ++-% %@@@@@@#  %=--------@**********@.   :-: .::
       :%----------#+-------*+ .:::   =#--#=.+###*-  +*--------+@##########%*     .-  :.
       #=----------=%=-------=*++===+*+==--**-.   .-**---------**++++++++===%.      ...
      +*------------+#---------------%@@@@%=-=+***+=-----------%=-----------+*
     .@**************%#=-------------=%@@#+-------------------=%------------=%
     +%***************#%=----------=+%+.##+=------------------**-------------*=
      :##*#%%%%@%%#***%=+*=--------++=#+#*++---------------=**@####**********#%
        +%@###@%##%%**%  .+****+==-----=---------------=+*+=. @************#%#*
        :@##%@%##%%%%%%=-.:#%##@@%%*-------------==+**#@=    :%******#***%*:
      =+%###@##%%+=*=-+*=%%###@%###@%*************+=---=**.  =%*#%%@%%%%@:
     #**%##%%##@=-+=-++-+@###@#####%%------------=++*#*=-=#+ =%%%###@%##%*
     @=%###@##%*--------%###@%##%%##@----------+%%%@%##%#==*%*@@%%%##%%##%*
     :+@###@##@+-------+@##%@#%%***@##--------=@####%@###@*+#=-++=*%##%%##@*+.
       @##%@##%%=------*%##@%#@===##=@=-------+%#####%@###%--*=-==-+@##@###%=%.
       ####@####%%#+=--#%##@##%##*%###--------*%######@%##%+--------%%#@%##@+#:
       :@##@%###%# :-==#%##@######@#=--------=#@####%#%@###%--------#%#%%##@*-
        .+*#%%%#=      :@##@%####%%=---------#*+%===#%#@###@=------+@%#%%##@:
                        =%%%@%##%#=----------+#=@++*%%#%%##@=--=+#%%###%%##@+
                         #++*##*=-------------+*@%%%###@%##@=++=##@####@##@%+
                         #=---------------------+@%####@##%+   #+-=%%##*+=:
                         ++----------------------=%%%%@@%%*   +*-=#-
                         -#------------------------=**++=#-   @=-*=
                          %=------------==---------------%:   %=-#-
                          -#------------=%+--------------%*+:.#*-++
                           #+------------+@+------------=%==+*==**.
                           .%=-----------*+*+-----------*=-++++=.
                            *+-----------#- %=----------%.
                            -#----------=@  :%----------@
                             %----------+*   *+---------@
                       .-====@+=--------%:    %=--------@   ...
                      **=----=++-------+*     =*--------%***+++**.
                     =*---------------=%.      %=---------------=%
                     -#-------------=**        +#==--------------@
                      =**++++++**++=-.          .:=++**++=====++#=
                         .::::.                        .::----:.
`

// https://fsymbols.com/generators/carty/
const banner = `
█▀ ▀█▀ ▄▀█ █▀▀ █▀▀ ▀█▀ █ █▀▄▀█ █▀▀ ▄▄ █ █▄░█ █▀▀ █▀█ █▀█ █▀▄▀█ █▀▀ █▀█
▄█ ░█░ █▀█ █▄█ ██▄ ░█░ █ █░▀░█ ██▄ ░░ █ █░▀█ █▀░ █▄█ █▀▄ █░▀░█ ██▄ █▀▄

`

func PrintBanner() {
	// Output banner + version output
	color.Magenta(logo)
	color.Cyan(banner)
	resp := goVersion.FuncWithOutput(false, version, commit, date, output)
	color.Magenta(resp + "\n")
}
