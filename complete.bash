function _temple()
{
    COMPREPLY=($(SHELL_AUTOCOMPLETE=1 ./temple "${COMP_LINE}"))
    [[ $COMPREPLY = */ ]] && compopt -o nospace     # don't add space after directories
}

complete -F _temple temple
