#!/bin/sh
# mock file that expects only these inputs:
# iptables -t nat -A PREROUTING -p tcp -s $sip --dport $dport -j DNAT --to-destination $dest
# iptables -t nat -D PREROUTING -p tcp -s $sip --dport $dport -j DNAT --to-destination $dest

if [ -z "${MITMGUI_TESTING_TESTDB}" ]; then
    DB=$(cd `dirname $0`; pwd)/test.db
else
    DB="${MITMGUI_TESTING_TESTDB}"
fi

if [ -z "${14}" ]; then
    echo "unexpected input"
    exit 2
fi

OP=$3
SIP=$8
DPORT=${10}
DEST=${14}

if [ "${OP}" == "-A" ]; then
    echo "${SIP}|${DPORT}|${DEST}" >> $DB
fi

if [ "${OP}" == "-D" ]; then
    # this should only delete first match (as iptables does) but it doesn't matter
    if ! cat $DB | grep ^"${SIP}|${DPORT}|${DEST}"$; then
        echo "iptables: No chain/target/match by that name." > /dev/stderr
        exit 1
    fi
    cat $DB | grep -v ^"${SIP}|${DPORT}|${DEST}"$ > $DB.tmp
    mv $DB.tmp $DB
fi
