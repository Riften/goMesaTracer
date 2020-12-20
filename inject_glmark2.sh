if [ -n "$GLM2_REPO_DIR" ]; then
	echo "Inject header to glmark2 at ${GLM2_REPO_DIR}"
	pathlist=("${GLM2_REPO_DIR}/src/")
	# shellcheck disable=SC2068
	for p in ${pathlist[@]}
	do
		echo "...Inject to ${p}"
		cp build/libMesaTracer.h ${p}
	done
else
	echo "No glmark2 repo specified."
	echo "Please specify glmark2 repo by \$GLM2_REPO_DIR."
fi
