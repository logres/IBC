# Generated by Django 3.1.14 on 2024-04-23 18:19

import api.models
import api.utils.common
from django.conf import settings
import django.contrib.auth.models
import django.contrib.postgres.fields
import django.core.validators
from django.db import migrations, models
import django.db.models.deletion
import django.utils.timezone


class Migration(migrations.Migration):

    initial = True

    dependencies = [
        ('auth', '0012_alter_user_first_name_max_length'),
    ]

    operations = [
        migrations.CreateModel(
            name='UserProfile',
            fields=[
                ('password', models.CharField(max_length=128, verbose_name='password')),
                ('last_login', models.DateTimeField(blank=True, null=True, verbose_name='last login')),
                ('is_superuser', models.BooleanField(default=False, help_text='Designates that this user has all permissions without explicitly assigning them.', verbose_name='superuser status')),
                ('first_name', models.CharField(blank=True, max_length=150, verbose_name='first name')),
                ('last_name', models.CharField(blank=True, max_length=150, verbose_name='last name')),
                ('is_staff', models.BooleanField(default=False, help_text='Designates whether the user can log into this admin site.', verbose_name='staff status')),
                ('is_active', models.BooleanField(default=True, help_text='Designates whether this user should be treated as active. Unselect this instead of deleting accounts.', verbose_name='active')),
                ('date_joined', models.DateTimeField(default=django.utils.timezone.now, verbose_name='date joined')),
                ('id', models.UUIDField(default=api.utils.common.make_uuid, help_text='ID of user', primary_key=True, serialize=False)),
                ('email', models.EmailField(db_index=True, max_length=254, unique=True)),
                ('username', models.CharField(default='', help_text='Name of user', max_length=64)),
                ('role', models.CharField(choices=[('admin', 'Admin'), ('operator', 'Operator'), ('user', 'User')], default=2, max_length=64)),
                ('groups', models.ManyToManyField(blank=True, help_text='The groups this user belongs to. A user will get all permissions granted to each of their groups.', related_name='user_set', related_query_name='user', to='auth.Group', verbose_name='groups')),
                ('user_permissions', models.ManyToManyField(blank=True, help_text='Specific permissions for this user.', related_name='user_set', related_query_name='user', to='auth.Permission', verbose_name='user permissions')),
            ],
            options={
                'verbose_name': 'User Info',
                'verbose_name_plural': 'User Info',
                'ordering': ['-date_joined'],
            },
            managers=[
                ('objects', django.contrib.auth.models.UserManager()),
            ],
        ),
        migrations.CreateModel(
            name='Agent',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, help_text='ID of agent', primary_key=True, serialize=False)),
                ('name', models.CharField(default='agent-b749b0c08b9a4dfc86b37489a6be7d3a', help_text='Agent name, can be generated automatically.', max_length=64)),
                ('urls', models.URLField(blank=True, help_text='Agent URL', null=True)),
                ('status', models.CharField(choices=[('inactive', 'Inactive'), ('active', 'Active')], default='active', help_text='Status of agent', max_length=10)),
                ('type', models.CharField(choices=[('docker', 'Docker'), ('kubernetes', 'Kubernetes')], default='docker', help_text='Type of agent', max_length=32)),
                ('config_file', models.FileField(blank=True, help_text='Config file for agent', max_length=256, upload_to=api.models.get_agent_config_file_path)),
                ('created_at', models.DateTimeField(auto_now_add=True, help_text='Create time of agent')),
                ('free_ports', django.contrib.postgres.fields.ArrayField(base_field=models.IntegerField(blank=True), help_text='Agent free ports.', null=True, size=None)),
            ],
            options={
                'ordering': ('-created_at',),
            },
        ),
        migrations.CreateModel(
            name='BPMN',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of userJoinOrgInvite', primary_key=True, serialize=False, unique=True)),
                ('name', models.CharField(blank=True, help_text='Name of Bpmn', max_length=255, null=True)),
                ('participants', models.TextField(blank=True, help_text='participants of BpmnStoragedFile', null=True)),
                ('bpmnContent', models.TextField(help_text='content of bpmn file')),
                ('svgContent', models.TextField(help_text='content of svg file')),
            ],
        ),
        migrations.CreateModel(
            name='Consortium',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of Consortium', primary_key=True, serialize=False, unique=True)),
                ('name', models.TextField(help_text='name of Consortium')),
            ],
        ),
        migrations.CreateModel(
            name='Environment',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of environment', primary_key=True, serialize=False, unique=True)),
                ('name', models.TextField(help_text='name of environment')),
                ('status', models.CharField(default='CREATED', help_text='status of environment,can be CREATED|INITIALIZED|STARTED|ACTIVATED|FIREFLY', max_length=32)),
                ('create_at', models.DateTimeField(auto_now_add=True, help_text='create time of environment')),
                ('consortium', models.ForeignKey(help_text='consortium of environment', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.consortium')),
            ],
        ),
        migrations.CreateModel(
            name='FabricPeer',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='', help_text='Name of peer node', max_length=64)),
                ('gossip_use_leader_reflection', models.BooleanField(default=True, help_text='Gossip use leader reflection')),
                ('gossip_org_leader', models.BooleanField(default=False, help_text='Gossip org leader')),
                ('gossip_skip_handshake', models.BooleanField(default=True, help_text='Gossip skip handshake')),
                ('local_msp_id', models.CharField(default='', help_text='Local msp id of peer node', max_length=64)),
            ],
        ),
        migrations.CreateModel(
            name='FabricResourceSet',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, help_text='ID of organization', primary_key=True, serialize=False)),
                ('name', models.CharField(default='', help_text='Name of organization', max_length=64)),
                ('created_at', models.DateTimeField(auto_now_add=True)),
                ('msp', models.TextField(help_text='msp of organization', null=True)),
                ('tls', models.TextField(help_text='tls of organization', null=True)),
                ('org_type', models.CharField(choices=[('userorg', 'USERORG'), ('systemorg', 'SYSTEMORG')], help_text='Organization type', max_length=32)),
            ],
            options={
                'ordering': ('-created_at',),
            },
        ),
        migrations.CreateModel(
            name='LoleidoMemebership',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of LoleidoMemebership', primary_key=True, serialize=False, unique=True)),
                ('role', models.CharField(choices=[('Member', 'Member'), ('Admin', 'Admin'), ('Owner', 'Owner')], default='Member', help_text='role of LoleidoMemebership', max_length=32)),
            ],
        ),
        migrations.CreateModel(
            name='LoleidoOrganization',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of LoleidoOrganization', primary_key=True, serialize=False, unique=True)),
                ('name', models.TextField(help_text='name of LoleidoOrganization')),
                ('members', models.ManyToManyField(help_text='related user_id', related_name='orgs', through='api.LoleidoMemebership', to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.CreateModel(
            name='Membership',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of membership', primary_key=True, serialize=False, unique=True)),
                ('name', models.TextField(help_text='name of membership')),
                ('create_at', models.DateTimeField(auto_now_add=True, help_text='create time of membership')),
                ('primary_contact_email', models.EmailField(help_text='primary contact email of membership', max_length=254, null=True)),
                ('consortium', models.ForeignKey(help_text='related consortium_id', on_delete=django.db.models.deletion.CASCADE, to='api.consortium')),
                ('loleido_organization', models.ForeignKey(help_text='related loleido_organization_id', on_delete=django.db.models.deletion.CASCADE, to='api.loleidoorganization')),
            ],
        ),
        migrations.CreateModel(
            name='Network',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, help_text='ID of network', primary_key=True, serialize=False)),
                ('name', models.CharField(default='netowrk-da433bbfa27b48e98f29751ade22feaf', help_text='network name, can be generated automatically.', max_length=64)),
                ('type', models.CharField(default='fabric', help_text="Type of network, ['fabric']", max_length=64)),
                ('version', models.CharField(default='', help_text="\n    Version of network.\n    Fabric supported versions: ['1.4.2', '1.5']\n    ", max_length=64)),
                ('created_at', models.DateTimeField(auto_now_add=True, help_text='Create time of network')),
                ('consensus', models.CharField(default='raft', help_text='Consensus of network', max_length=128)),
                ('genesisblock', models.TextField(help_text='genesis block', null=True)),
                ('database', models.CharField(default='leveldb', help_text='database of network', max_length=128)),
            ],
            options={
                'ordering': ('-created_at',),
            },
        ),
        migrations.CreateModel(
            name='Node',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, help_text='ID of node', primary_key=True, serialize=False)),
                ('name', models.CharField(default='', help_text='Node name', max_length=64)),
                ('type', models.CharField(help_text="\n    Node type defined for network.\n    Fabric available types: ['ca', 'orderer', 'peer']\n    ", max_length=64)),
                ('urls', models.JSONField(blank=True, default=dict, help_text='URL configurations for node', null=True)),
                ('created_at', models.DateTimeField(auto_now_add=True, help_text='Create time of network')),
                ('status', models.CharField(choices=[('created', 'Created'), ('restarting', 'Restarting'), ('running', 'Running'), ('removing', 'Removing'), ('paused', 'Paused'), ('exited', 'Exited'), ('dead', 'Dead')], default='created', help_text='Status of node', max_length=64)),
                ('config_file', models.TextField(help_text='Config file of node', null=True)),
                ('msp', models.TextField(help_text='msp of node', null=True)),
                ('tls', models.TextField(help_text='tls of node', null=True)),
                ('cid', models.CharField(default='', help_text='id used in agent, such as container id', max_length=256)),
                ('agent', models.ForeignKey(help_text='Agent of node', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='node', to='api.agent')),
                ('fabric_resource_set', models.ForeignKey(help_text='Organization of node', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='node', to='api.fabricresourceset')),
                ('user', models.ForeignKey(help_text='User of node', null=True, on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
            options={
                'ordering': ('-created_at',),
            },
        ),
        migrations.CreateModel(
            name='NodeUser',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(default='', help_text='User name of node', max_length=64)),
                ('secret', models.CharField(default='', help_text='User secret of node', max_length=64)),
                ('user_type', models.CharField(choices=[('peer', 'Peer'), ('orderer', 'Orderer'), ('user', 'User')], default='peer', help_text='User type of node', max_length=64)),
                ('status', models.CharField(choices=[('registering', 'Registering'), ('registered', 'Registered'), ('fail', 'Fail')], default='registering', help_text='Status of node user', max_length=32)),
                ('attrs', models.CharField(default='', help_text='Attributes of node user', max_length=512)),
                ('node', models.ForeignKey(help_text='Node of user', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.node')),
            ],
            options={
                'ordering': ('id',),
            },
        ),
        migrations.CreateModel(
            name='PeerCa',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('address', models.CharField(default='', help_text='Node Address of ca', max_length=128)),
                ('certificate', models.FileField(blank=True, help_text='Certificate file for ca node.', max_length=256, null=True, upload_to=api.models.get_ca_certificate_path)),
                ('type', models.CharField(choices=[('tls', 'TLS'), ('signature', 'Signature'), ('both', 'Both')], default='signature', help_text='Type of ca node for peer', max_length=64)),
                ('node', models.ForeignKey(help_text='CA node of peer', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.node')),
                ('peer', models.ForeignKey(help_text='Peer node', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.fabricpeer')),
            ],
        ),
        migrations.CreateModel(
            name='ResourceSet',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of midOrg', primary_key=True, serialize=False, unique=True)),
                ('name', models.TextField(help_text='name of ResourceSet')),
                ('agent', models.ForeignKey(help_text='related agent_id', null=True, on_delete=django.db.models.deletion.SET_NULL, to='api.agent')),
                ('environment', models.ForeignKey(help_text='related environment_Id', on_delete=django.db.models.deletion.CASCADE, related_name='resource_sets', to='api.environment')),
                ('membership', models.ForeignKey(help_text='related membership_id', on_delete=django.db.models.deletion.CASCADE, to='api.membership')),
            ],
        ),
        migrations.CreateModel(
            name='UserPreference',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of userPreference', primary_key=True, serialize=False, unique=True)),
                ('last_active_consortium', models.ForeignKey(help_text='related consortium_id', null=True, on_delete=django.db.models.deletion.SET_NULL, to='api.consortium')),
                ('last_active_environment', models.ForeignKey(help_text='related environment_id', null=True, on_delete=django.db.models.deletion.SET_NULL, to='api.environment')),
                ('last_active_organization', models.ForeignKey(help_text='related middle_org_id', null=True, on_delete=django.db.models.deletion.SET_NULL, to='api.resourceset')),
                ('user', models.ForeignKey(help_text='related user_id', on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.CreateModel(
            name='UserJoinOrgInvitation',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of userJoinOrgInvite', primary_key=True, serialize=False, unique=True)),
                ('role', models.TextField(default='Member', help_text='role of userJoinOrgInvite')),
                ('message', models.TextField(help_text='message of userJoinOrgInvite')),
                ('status', models.CharField(choices=[('pending', 'Pending'), ('accepted', 'Accepted'), ('rejected', 'Rejected')], default='pending', help_text='status of userJoinOrgInvite', max_length=32)),
                ('create_at', models.DateTimeField(auto_now_add=True, help_text='Create time of userJoinOrgInvite')),
                ('invitor', models.ForeignKey(default=None, help_text='related user_id', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='sended_invitations', to=settings.AUTH_USER_MODEL)),
                ('loleido_organization', models.ForeignKey(help_text='related loleido_organization_id', on_delete=django.db.models.deletion.CASCADE, to='api.loleidoorganization')),
                ('user', models.ForeignKey(help_text='related user_id', on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.CreateModel(
            name='Port',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('external', models.IntegerField(default=0, help_text='External port', validators=[django.core.validators.MinValueValidator(1), django.core.validators.MaxValueValidator(65535)])),
                ('internal', models.IntegerField(default=0, help_text='Internal port', validators=[django.core.validators.MinValueValidator(1), django.core.validators.MaxValueValidator(65535)])),
                ('node', models.ForeignKey(help_text='Node of port', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='port', to='api.node')),
            ],
            options={
                'ordering': ('external',),
            },
        ),
        migrations.CreateModel(
            name='PeerCaUser',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('username', models.CharField(default='', help_text='If user not set, set username/password', max_length=64)),
                ('password', models.CharField(default='', help_text='If user not set, set username/password', max_length=64)),
                ('type', models.CharField(choices=[('peer', 'Peer'), ('orderer', 'Orderer'), ('user', 'User')], default='user', help_text='User type of ca', max_length=64)),
                ('peer_ca', models.ForeignKey(help_text='Peer Ca configuration', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.peerca')),
                ('user', models.ForeignKey(help_text='User of ca node', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.nodeuser')),
            ],
        ),
        migrations.CreateModel(
            name='LoleidoOrgJoinConsortiumInvitation',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of loleidoOrgJoinConsortiumInvite', primary_key=True, serialize=False, unique=True)),
                ('role', models.TextField(default='Member', help_text='role of loleidoOrgJoinConsortiumInvite')),
                ('message', models.TextField(help_text='message of loleidoOrgJoinConsortiumInvite')),
                ('status', models.CharField(choices=[('pending', 'Pending'), ('accepted', 'Accepted'), ('rejected', 'Rejected')], default='pending', help_text='status of LoleidoOrgJoinConsortiumInvite', max_length=32)),
                ('create_at', models.DateTimeField(auto_now_add=True, help_text='Create time of loleidoOrgJoinConsortiumInvite')),
                ('consortium', models.ForeignKey(help_text='related consortium_id', on_delete=django.db.models.deletion.CASCADE, to='api.consortium')),
                ('invitor', models.ForeignKey(default=None, help_text='related loleido_organization_id', on_delete=django.db.models.deletion.CASCADE, related_name='sended_invitations', to='api.loleidoorganization')),
                ('loleido_organization', models.ForeignKey(help_text='related loleido_organization_id', on_delete=django.db.models.deletion.CASCADE, to='api.loleidoorganization')),
            ],
        ),
        migrations.AddField(
            model_name='loleidomemebership',
            name='loleido_organization',
            field=models.ForeignKey(help_text='related loleido_organization_id', on_delete=django.db.models.deletion.CASCADE, to='api.loleidoorganization'),
        ),
        migrations.AddField(
            model_name='loleidomemebership',
            name='user',
            field=models.ForeignKey(help_text='related user_id', on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL),
        ),
        migrations.CreateModel(
            name='KubernetesConfig',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('credential_type', models.CharField(choices=[('cert_key', 'CertKey'), ('config', 'Config'), ('username_password', 'UsernamePassword')], default='cert_key', help_text='Credential type of k8s', max_length=32)),
                ('enable_ssl', models.BooleanField(default=False, help_text='Whether enable ssl for api')),
                ('ssl_ca', models.TextField(blank=True, default='', help_text='Ca file content for ssl')),
                ('nfs_server', models.CharField(blank=True, default='', help_text='NFS server address for k8s', max_length=256)),
                ('parameters', models.JSONField(blank=True, default=dict, help_text='Extra parameters for kubernetes', null=True)),
                ('cert', models.TextField(blank=True, default='', help_text='Cert content for k8s')),
                ('key', models.TextField(blank=True, default='', help_text='Key content for k8s')),
                ('username', models.CharField(blank=True, default='', help_text='Username for k8s credential', max_length=128)),
                ('password', models.CharField(blank=True, default='', help_text='Password for k8s credential', max_length=128)),
                ('agent', models.ForeignKey(help_text='Agent of kubernetes config', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.agent')),
            ],
        ),
        migrations.CreateModel(
            name='Firefly',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of firefly', primary_key=True, serialize=False, unique=True)),
                ('org_name', models.CharField(help_text='org name of firefly', max_length=128)),
                ('core_url', models.TextField(help_text='name of core url')),
                ('sandbox_url', models.TextField(help_text='name of sandbox url')),
                ('fab_connect_url', models.TextField(blank=True, help_text='name of fabconnect url', null=True)),
                ('resource_set', models.ForeignKey(help_text='related resource_set id', on_delete=django.db.models.deletion.CASCADE, to='api.resourceset')),
            ],
        ),
        migrations.CreateModel(
            name='File',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, help_text='ID of file', primary_key=True, serialize=False)),
                ('name', models.CharField(default='', help_text='File name', max_length=64)),
                ('file', models.FileField(blank=True, help_text='File', max_length=256, upload_to=api.models.get_file_path)),
                ('created_at', models.DateTimeField(auto_now_add=True, help_text='Create time of agent')),
                ('type', models.CharField(choices=[('certificate', 'Certificate')], default='certificate', help_text='File type', max_length=32)),
                ('organization', models.ForeignKey(help_text='Organization of file', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.fabricresourceset')),
            ],
            options={
                'ordering': ('-created_at',),
            },
        ),
        migrations.AddField(
            model_name='fabricresourceset',
            name='network',
            field=models.ForeignKey(help_text='Network to which the organization belongs', null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='organization', to='api.network'),
        ),
        migrations.AddField(
            model_name='fabricresourceset',
            name='resource_set',
            field=models.ForeignKey(help_text='Resource set to which the fabric resourceset belongs', null=True, on_delete=django.db.models.deletion.SET_NULL, related_name='sub_resource_set', to='api.resourceset'),
        ),
        migrations.CreateModel(
            name='FabricIdentity',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of FabricIdentity', primary_key=True, serialize=False, unique=True)),
                ('name', models.TextField(help_text='name of FabricIdentity')),
                ('signer', models.TextField(blank=True, help_text='signer of FabricIdentity', null=True)),
                ('secret', models.TextField(blank=True, help_text='secret of FabricIdentity', null=True)),
                ('create_at', models.DateTimeField(auto_now_add=True, help_text='Create time of FabricIdentity')),
                ('environment', models.ForeignKey(help_text='related environment_id', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.environment')),
                ('membership', models.ForeignKey(help_text='related membership_id', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.membership')),
            ],
        ),
        migrations.CreateModel(
            name='FabricCA',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('admin_name', models.CharField(default='admin', help_text='Admin username for ca server', max_length=32)),
                ('admin_password', models.CharField(default='adminpw', help_text='Admin password for ca server', max_length=32)),
                ('hosts', models.JSONField(blank=True, default=list, help_text='Hosts for ca', null=True)),
                ('type', models.CharField(choices=[('tls', 'TLS'), ('signature', 'Signature'), ('both', 'Both')], default='signature', help_text='Fabric ca server type', max_length=32)),
                ('node', models.ForeignKey(help_text='Node of ca', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.node')),
            ],
        ),
        migrations.AddField(
            model_name='environment',
            name='network',
            field=models.ForeignKey(null=True, on_delete=django.db.models.deletion.DO_NOTHING, to='api.network'),
        ),
        migrations.AddField(
            model_name='consortium',
            name='orgs',
            field=models.ManyToManyField(help_text='related loleido_organization_id', related_name='consortiums', through='api.Membership', to='api.LoleidoOrganization'),
        ),
        migrations.CreateModel(
            name='Channel',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of Channel', primary_key=True, serialize=False, unique=True)),
                ('name', models.CharField(help_text='name of channel', max_length=128)),
                ('create_ts', models.DateTimeField(auto_now_add=True, help_text='Create time of Channel')),
                ('config', models.JSONField(blank=True, default=dict, help_text='Channel config', null=True)),
                ('fabric_resource_set', models.ManyToManyField(help_text='the organization of the channel', related_name='channels', to='api.FabricResourceSet')),
                ('network', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='api.network')),
                ('orderers', models.ManyToManyField(help_text='Orderer list in the channel', to='api.Node')),
            ],
        ),
        migrations.CreateModel(
            name='ChainCode',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of ChainCode', primary_key=True, serialize=False, unique=True)),
                ('name', models.CharField(help_text='name of chainCode', max_length=128)),
                ('version', models.CharField(help_text='version of chainCode', max_length=128)),
                ('language', models.CharField(help_text='language of chainCode', max_length=128)),
                ('create_ts', models.DateTimeField(auto_now_add=True, help_text='Create time of chainCode')),
                ('creator', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to='api.loleidoorganization')),
                ('environment', models.ForeignKey(help_text='environment of chainCode', on_delete=django.db.models.deletion.CASCADE, to='api.environment')),
            ],
        ),
        migrations.CreateModel(
            name='BPMNInstance',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of BPMNInstance', primary_key=True, serialize=False, unique=True)),
                ('name', models.CharField(blank=True, help_text='Name of BPMNInstance', max_length=255, null=True)),
                ('status', models.CharField(choices=[('Initiated', 'Initiated'), ('Fullfilled', 'Fullfilled'), ('Generated', 'Generated'), ('Installed', 'Installed'), ('Registered', 'Registered')], default='pending', help_text='status of BPMNInstance', max_length=32)),
                ('create_at', models.DateTimeField(auto_now_add=True, help_text='Create time of BPMNInstance')),
                ('update_at', models.DateTimeField(auto_now=True, help_text='Update time of BPMNInstance')),
                ('chaincode_content', models.TextField(blank=True, default=None, help_text='content of chaincode file', null=True)),
                ('firefly_url', models.TextField(blank=True, help_text='firefly url of BPMNInstance', null=True)),
                ('ffiContent', models.TextField(blank=True, default=None, help_text='content of ffi file', null=True)),
                ('bpmn', models.ForeignKey(help_text='related bpmn_id', on_delete=django.db.models.deletion.CASCADE, to='api.bpmn')),
                ('chaincode', models.ForeignKey(help_text='related chaincode_id', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.chaincode')),
                ('environment', models.ForeignKey(help_text='related environment_id', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.environment')),
            ],
        ),
        migrations.CreateModel(
            name='BPMNBindingRecord',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of BPMNBindingRecord', primary_key=True, serialize=False, unique=True)),
                ('participant_id', models.CharField(blank=True, help_text='ID of participant', max_length=255, null=True)),
                ('bpmn_instance', models.ForeignKey(help_text='related bpmn_instance_id', on_delete=django.db.models.deletion.CASCADE, to='api.bpmninstance')),
                ('membership', models.ForeignKey(help_text='related membership_id', on_delete=django.db.models.deletion.CASCADE, to='api.membership')),
            ],
        ),
        migrations.AddField(
            model_name='bpmn',
            name='consortium',
            field=models.ForeignKey(help_text='related consortium_id', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.consortium'),
        ),
        migrations.AddField(
            model_name='bpmn',
            name='organization',
            field=models.ForeignKey(help_text='related organization_id', on_delete=django.db.models.deletion.CASCADE, to='api.loleidoorganization'),
        ),
        migrations.CreateModel(
            name='APISecretKey',
            fields=[
                ('id', models.UUIDField(default=api.utils.common.make_uuid, editable=False, help_text='ID of APISecretKey', primary_key=True, serialize=False, unique=True)),
                ('key', models.CharField(blank=True, help_text='key of APISecretKey', max_length=255, null=True)),
                ('key_secret', models.CharField(blank=True, help_text='key_secret of APISecretKey', max_length=255, null=True)),
                ('create_at', models.DateTimeField(auto_now_add=True, help_text='Create time of APISecretKey')),
                ('environment', models.ForeignKey(help_text='related environment_id', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.environment')),
                ('membership', models.ForeignKey(help_text='related membership_id', null=True, on_delete=django.db.models.deletion.CASCADE, to='api.membership')),
                ('user', models.ForeignKey(help_text='related user_id', on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.AddField(
            model_name='agent',
            name='organization',
            field=models.ForeignKey(help_text='Organization of agent', null=True, on_delete=django.db.models.deletion.CASCADE, related_name='agents', to='api.loleidoorganization'),
        ),
        migrations.AlterUniqueTogether(
            name='loleidomemebership',
            unique_together={('loleido_organization', 'user')},
        ),
    ]
