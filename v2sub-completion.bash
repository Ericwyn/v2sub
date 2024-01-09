#!/bin/bash

_v2sub_completion() {
    local cur prev opts sub_opts ser_opts conf_opts conn_opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="-sub -ser -conf -conn -v -h"
    sub_opts="add update customer remove list"
    ser_opts="list set setx setflush speedtest"
    conf_opts="sport hport lconn list"
    conn_opts="start start-pac kill"

    case "${prev}" in
        -sub)
            COMPREPLY=( $(compgen -W "${sub_opts}" -- ${cur}) )
            return 0
            ;;
        -ser)
            COMPREPLY=( $(compgen -W "${ser_opts}" -- ${cur}) )
            return 0
            ;;
        -conf)
            COMPREPLY=( $(compgen -W "${conf_opts}" -- ${cur}) )
            return 0
            ;;
        -conn)
            COMPREPLY=( $(compgen -W "${conn_opts}" -- ${cur}) )
            return 0
            ;;
        *)
            COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
            return 0
            ;;
    esac
}
complete -F _v2sub_completion v2sub
