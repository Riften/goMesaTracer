if [ -n "$MESA_REPO_DIR" ]; then
	echo "Inject header to mesa at ${MESA_REPO_DIR}"
	pathlist=("${MESA_REPO_DIR}/src/mesa/main/")
	for p in ${pathlist[@]}
	do
		echo "...Inject to ${p}"
		cp build/libMesaTracer.h ${p}
	done
else
	echo "No mesa repo specified."
	echo "Please specify mesa repo by \$MESA_REPO_DIR."
fi
