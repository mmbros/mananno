#!/bin/sh
### BEGIN INIT INFO
# Provides:          mananno
# Required-Start:    $local_fs $network $named $time $syslog
# Required-Stop:     $local_fs $network $named $time $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: mananno init script
# Description:       MMbros personal site
### END INIT INFO

# Author: MMbros
# Date: 2016-02-28

# Using the lsb functions to perform the operations.
. /lib/lsb/init-functions

# Process name ( For display )
NAME=mananno

# Daemon working dir
DIR=/home/pi/mananno

# Daemon name, where is the actual executable
DAEMON=$DIR/$NAME

# Daemon user
USERNAME=pi

# pid file for the daemon
PIDFILE=/var/run/$NAME.pid


# If the daemon is not there, then exit.
test -x $DAEMON || exit 5


daemon_start() {
  # Checked the PID file exists and check the actual status of process
  if [ -e $PIDFILE ]; then
    status_of_proc -p $PIDFILE $DAEMON "$NAME process" && status="0" || status="$?"
    # If the status is SUCCESS then don't need to start again.
    if [ $status = "0" ]; then
      exit # Exit
    fi
  fi
  # Start the daemon.
  log_daemon_msg "Starting the process" "$NAME"
  # Start the daemon with the help of start-stop-daemon
  # Log the message appropriately
  if start-stop-daemon --start --chuid $USERNAME --chdir $DIR --pidfile $PIDFILE --make-pidfile --background --exec $DAEMON ; then
    log_end_msg 0
  else
    log_end_msg 1
  fi
}


daemon_stop() {
  # Stop the daemon.
  if [ -e $PIDFILE ]; then
    status_of_proc -p $PIDFILE $DAEMON "Stoppping the $NAME process" && status="0" || status="$?"
    if [ "$status" = 0 ]; then
      start-stop-daemon --stop --quiet --oknodo --pidfile $PIDFILE
      /bin/rm -rf $PIDFILE
    fi
  else
    log_daemon_msg "$NAME process is not running"
    log_end_msg 0
  fi
}

daemon_status() {
  # Check the status of the process.
  if [ -e $PIDFILE ]; then
    status_of_proc -p $PIDFILE $DAEMON "$NAME process" && exit 0 || exit $?
  else
    log_daemon_msg "$NAME Process is not running"
    log_end_msg 0
  fi
}

daemon_reload() {
  # Reload the process. Basically sending some signal to a daemon to reload
  # its configurations.
  if [ -e $PIDFILE ]; then
    start-stop-daemon --stop --signal USR1 --quiet --pidfile $PIDFILE --name $NAME
    log_success_msg "$NAME process reloaded successfully"
  else
    log_failure_msg "$PIDFILE does not exists"
  fi
}

daemon_uninstall() {
  echo -n "Are you really sure you want to uninstall \"$NAME\" service?\nThat cannot be undone. [yes|No] "
  local SURE
  read SURE
  if [ "$SURE" = "yes" ]; then
    daemon_stop
    # echo "Notice: log file is not be removed: '$LOGFILE'" >&2
    update-rc.d -f $NAME remove
    rm -fv "$0"
  fi
}


case $1 in
  start)
    daemon_start
    ;;
  stop)
    daemon_stop
    ;;
  restart)
    # Restart the daemon.
    $0 stop && sleep 2 && $0 start
    ;;
  status)
    daemon_status
    ;;
  reload)
    daemon_reload
    ;;
  uninstall)
    daemon_uninstall
    ;;
  *)
    # For invalid arguments, print the usage message.
	echo "Usage: $(basename $0) {start|stop|restart|reload|status|uninstall}"
    exit 2
    ;;
esac
