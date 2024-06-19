# cloudinfo

Get information on the currently running cloud no matter the provider

## Adding a cloud provider

Submit the output of the following command in a [new issue](https://github.com/KarpelesLab/cloudinfo/issues/new):

    for foo in /sys/class/dmi/id/*; do echo -n "$(basename $foo) = "; cat $foo; done

