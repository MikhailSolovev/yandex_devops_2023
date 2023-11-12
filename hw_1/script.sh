ls -l | sed 's/\(^[dl-]\)/ \1/'

ls -l | awk 'match($0, /^[dl-]/) { $0 = substr($0, 1, RSTART - 1) " " substr($0, RSTART) } 1'
