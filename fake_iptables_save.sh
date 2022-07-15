#!/bin/sh

if [ -z "${MITMGUI_TESTING_TESTDB}" ]; then
    DB=$(cd `dirname $0`; pwd)/test.db
else
    DB="${MITMGUI_TESTING_TESTDB}"
fi

touch $DB

for line in $(cat $DB); do
    SIP=$(echo $line | cut -d '|' -f1)
    DPORT=$(echo $line | cut -d '|' -f2)
    DEST=$(echo $line | cut -d '|' -f3)
    echo "-A PREROUTING -s $SIP/32 -p tcp -m tcp --dport $DPORT -j DNAT --to-destination $DEST"
done
