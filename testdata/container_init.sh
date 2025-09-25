#!/bin/sh

sed -i "s/^#jit = on/jit = off/g" /var/lib/postgresql/data/postgresql.conf
current_datetime_iso8601=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE EXTENSION IF NOT EXISTS age;
    LOAD 'age';
    SET search_path = ag_catalog, $POSTGRES_USER, public;
    ALTER ROLE $POSTGRES_USER SET search_path TO ag_catalog,"\$user",public;

    SELECT ag_catalog.create_graph('cdim_graph');

    SELECT CREATE_VLABEL('cdim_graph', 'CXLswitch');
    SELECT CREATE_VLABEL('cdim_graph', 'Annotation');
    SELECT CREATE_VLABEL('cdim_graph', 'CPU');
    SELECT CREATE_VLABEL('cdim_graph', 'Memory');
    SELECT CREATE_VLABEL('cdim_graph', 'Storage');
    SELECT CREATE_VLABEL('cdim_graph', 'NetworkInterface');
    SELECT CREATE_VLABEL('cdim_graph', 'GraphicController');
    SELECT CREATE_VLABEL('cdim_graph', 'VirtualMedia');
    SELECT CREATE_VLABEL('cdim_graph', 'Accelerator');
    SELECT CREATE_VLABEL('cdim_graph', 'DSP');
    SELECT CREATE_VLABEL('cdim_graph', 'FPGA');
    SELECT CREATE_VLABEL('cdim_graph', 'GPU');
    SELECT CREATE_VLABEL('cdim_graph', 'UnknownProcessor');
    SELECT CREATE_VLABEL('cdim_graph', 'Node');
    SELECT CREATE_VLABEL('cdim_graph', 'Rack');
    SELECT CREATE_VLABEL('cdim_graph', 'Chassis');
    SELECT CREATE_VLABEL('cdim_graph', 'NotDetectedDevice');
    SELECT CREATE_VLABEL('cdim_graph', 'ResourceGroups');

    SELECT CREATE_ELABEL('cdim_graph', 'Connect');
    SELECT CREATE_ELABEL('cdim_graph', 'Compose');
    SELECT CREATE_ELABEL('cdim_graph', 'Mount');
    SELECT CREATE_ELABEL('cdim_graph', 'Attach');
    SELECT CREATE_ELABEL('cdim_graph', 'Have');
    SELECT CREATE_ELABEL('cdim_graph', 'Include');
    SELECT CREATE_ELABEL('cdim_graph', 'NotDetected');

    CREATE INDEX cdim_graph_CXLswitch_idx ON cdim_graph."CXLswitch" USING gin (properties);
    CREATE INDEX cdim_graph_Annotation_idx ON cdim_graph."Annotation" USING gin (properties);
    CREATE INDEX cdim_graph_CPU_idx ON cdim_graph."CPU" USING gin (properties);
    CREATE INDEX cdim_graph_Memory_idx ON cdim_graph."Memory" USING gin (properties);
    CREATE INDEX cdim_graph_Storage_idx ON cdim_graph."Storage" USING gin (properties);
    CREATE INDEX cdim_graph_NetworkInterface_idx ON cdim_graph."NetworkInterface" USING gin (properties);
    CREATE INDEX cdim_graph_GraphicController_idx ON cdim_graph."GraphicController" USING gin (properties);
    CREATE INDEX cdim_graph_VirtualMedia_idx ON cdim_graph."VirtualMedia" USING gin (properties);
    CREATE INDEX cdim_graph_Accelerator_idx ON cdim_graph."Accelerator" USING gin (properties);
    CREATE INDEX cdim_graph_DSP_idx ON cdim_graph."DSP" USING gin (properties);
    CREATE INDEX cdim_graph_FPGA_idx ON cdim_graph."FPGA" USING gin (properties);
    CREATE INDEX cdim_graph_GPU_idx ON cdim_graph."GPU" USING gin (properties);
    CREATE INDEX cdim_graph_UnknownProcessor_idx ON cdim_graph."UnknownProcessor" USING gin (properties);
    CREATE INDEX cdim_graph_Node_idx ON cdim_graph."Node" USING gin (properties);
    CREATE INDEX cdim_graph_Rack_idx ON cdim_graph."Rack" USING gin (properties);
    CREATE INDEX cdim_graph_Chassis_idx ON cdim_graph."Chassis" USING gin (properties);
    CREATE INDEX cdim_graph_NotDetectedDevice_idx ON cdim_graph."NotDetectedDevice" USING gin (properties);
    CREATE INDEX cdim_graph_ResourceGroups_idx ON cdim_graph."ResourceGroups" USING gin (properties);

    SELECT * FROM cypher('cdim_graph', \$\$ CREATE (a: NotDetectedDevice) \$\$) AS (a agtype);

    SELECT * FROM cypher('cdim_graph', \$\$ CREATE (a: ResourceGroups {id: "00000000-0000-7000-8000-000000000000", name: "default", description: "default group", createdAt: "$current_datetime_iso8601", updatedAt: "$current_datetime_iso8601"}) \$\$) AS (a agtype);
EOSQL
