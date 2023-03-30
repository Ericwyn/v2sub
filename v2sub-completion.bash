#/usr/bin/env bash
function _myscript(){
    if [[ "${COMP_CWORD}" == "1" ]];then
        COMP_WORD="-sub -ser -conf -conn -v -h"
        COMPREPLY=($(compgen -W "$COMP_WORD" -- ${COMP_WORDS[${COMP_CWORD}]}))
    else
        case ${COMP_WORDS[$[$COMP_CWORD-1]]} in
        -sub)
        COMP_WORD_2="add update remove list"
        COMPREPLY=($(compgen -W "${COMP_WORD_2}" ${COMP_WORDS[${COMP_CWORD}]}))
        ;;

        -ser)
        COMP_WORD_2="list set setx setflush speedtest"
        COMPREPLY=($(compgen -W "${COMP_WORD_2}" ${COMP_WORDS[${COMP_CWORD}]}))
        ;;

        -conf)
        COMP_WORD_2="sport hport lconn list"
        COMPREPLY=($(compgen -W "${COMP_WORD_2}" ${COMP_WORDS[${COMP_CWORD}]}))
        ;;

        -conn)
        COMP_WORD_2="start start-pac kill"
        COMPREPLY=($(compgen -W "${COMP_WORD_2}" ${COMP_WORDS[${COMP_CWORD}]}))
        ;;

        esac
    fi
}
# 注册命令补全函数
complete -F _myscript v2sub
