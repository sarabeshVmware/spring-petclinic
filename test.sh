for dir in /Users/mdogra/tap-packages/packages/*; do
    
    for file in $dir/*; do
        if [ -d "$file" ]; then
 	    if [[ $file == *"metadata.yaml"* ]]; then
			  echo "Skipping metadata file"
	    else
                file_path=$file/package.yaml
		new_file=$file.yaml 
		# echo "mv $file_path $new_file"
		# mv $file_path $new_file
		echo "rmdir $file"
		rmdir $file
	    fi
        fi
    done
done

