<?php

namespace Database\Seeders;

use App\Models\StaticSiteContent;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class StaticSiteContentSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('static_site_contents')->delete();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/staticSiteContents.csv'), ';');

        $model = new StaticSiteContent();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
