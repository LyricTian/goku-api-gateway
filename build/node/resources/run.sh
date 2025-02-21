
#!/bin/bash
### BEGIN INIT INFO
#
# Provides:	 location_server
# Required-Start:	$local_fs  $remote_fs
# Required-Stop:	$local_fs  $remote_fs
# Default-Start: 	2 3 4 5
# Default-Stop: 	0 1 6
# Short-Description:	initscript
# Description: 	This file should be used to construct scripts to be placed in /etc/init.d.
#
### END INIT INFO
 
## Fill in name of program here.
cd  $(dirname $0) # 当前位置跳到脚本位置
PROG="node"
PROG_PATH="$(pwd)" ## Not need, but sometimes helpful (if $PROG resides in /opt for example).

WORK_PATH="$PROG_PATH/work"

start() {
    if [[ -e "$WORK_PATH/$PROG.env" ]]; then
        source $WORK_PATH/$PROG.env
    fi

    ADMIN=$1
    PORT=$2

    if [[ "$PORT" = "" ]] ;then
        PORT=${ENV_PORT}
    fi
    if [[ "$PORT" = "" ]] ;then
        PORT="6689"
    fi
    if [[ "$ADMIN" = "" ]] ; then
        ADMIN=${ENV_ADMIN}
    fi

    if [[ "$ADMIN" = "" ]] ; then
        echo "start fail :need admin url"
        exit 1
    fi

    mkdir -p $WORK_PATH/logs
    echo -e "ENV_PORT=$PORT\nENV_ADMIN=$ADMIN" > $WORK_PATH/$PROG.env

    if [[ -e "$WORK_PATH/$PROG.pid" ]]; then
        ## Program is running, exit with error.
        echo "Error! $PROG is currently running!" 1>&2
        exit 1
    else
        time=$(date "+%Y%m%d-%H%M%S")
        ## Change from /dev/null to something like /var/log/$PROG if you want to save output.
        nohup $PROG_PATH/$PROG --port=$PORT --admin=$ADMIN 2>&1 >"$WORK_PATH/logs/stdout-$PROG-$time.log" &  pid=$!

        echo "$PROG started"
        echo $pid > "$WORK_PATH/$PROG.pid"
    fi
}
 
stop() {
    echo "begin stop"
    if [[ -e "$WORK_PATH/$PROG.pid" ]]; then
        ## Program is running, so stop it
        pid="$(cat $WORK_PATH/$PROG.pid)"
        if [[ "ps ax|grep $pid|grep '$PROG' |awk '{print \$1}'" != ""  ]];then
            kill $pid
            if [[ $? != 0 ]];then
                echo "$PROG stop error"
                exit 1
            fi
            rm -f  "$WORK_PATH/$PROG.pid"
            echo "$PROG stopped"
        fi
    else
        ## Program is not running, exit with error.
        echo "Note! $PROG not started!" 1>&2

    fi
}
 
## Check to see if we are running as root first.
## Found at http://www.cyberciti.biz/tips/shell-root-user-check-script.html
#if [[ "$(id -u)" != "0" ]]; then
#    echo "This script must be run as root" 1>&2
#    exit 1
#fi
#
case "$1" in
    start)
        start $2 $3
        exit 0
    ;;
    stop)
        stop
        exit 0
    ;;
    reload|restart|force-reload)
        stop
        start $2 $3
        exit 0
    ;;
    **)
        echo "Usage: $0 {start|stop|reload|restart|force-reload} [admin url] [port] " 1>&2
        exit 1
    ;;
esac
