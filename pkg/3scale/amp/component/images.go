package component

func ApicastImageURL() string {
	return "quay.io/3scale/apicast:latest"
}

func BackendImageURL() string {
	return "quay.io/3scale/apisonator:latest"
}

func SystemImageURL() string {
	return "quay.io/3scale/porta:latest"
}

func SystemSearchdImageURL() string {
	return "quay.io/3scale/searchd:latest"
}

func ZyncImageURL() string {
	return "quay.io/3scale/zync:latest"
}

func BackendRedisImageURL() string {
	return "quay.io/fedora/redis-6"
}

func SystemRedisImageURL() string {
	return "quay.io/fedora/redis-6"
}

func SystemMySQLImageURL() string {
	return "quay.io/sclorg/mysql-80-c8s"
}

func SystemPostgreSQLImageURL() string {
	return "quay.io/sclorg/postgresql-10-c8s"
}

func SystemMemcachedImageURL() string {
	return "mirror.gcr.io/library/memcached:1.5"
}

func ZyncPostgreSQLImageURL() string {
	return "quay.io/sclorg/postgresql-10-c8s"
}

func OCCLIImageURL() string {
	return "quay.io/openshift/origin-cli:4.7"
}
