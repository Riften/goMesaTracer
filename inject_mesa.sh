if [ -n "$MESA_REPO_DIR" ]; then
	echo "Inject header to mesa at ${MESA_REPO_DIR}"
	pathlist=("${MESA_REPO_DIR}/src/mesa/main/" "${MESA_REPO_DIR}/src/mesa/state_tracker" "${MESA_REPO_DIR}/src/mesa/drivers/x11" "${MESA_REPO_DIR}/src/glx" "${MESA_REPO_DIR}/src/gallium/auxiliary/util" "${MESA_REPO_DIR}/src/gallium/frontends/dri")
	for p in ${pathlist[@]}
	do
		echo "...Inject to ${p}"
		cp build/libMesaTracer.h ${p}
	done
else
	echo "No mesa repo specified."
	echo "Please specify mesa repo by \$MESA_REPO_DIR."
fi
