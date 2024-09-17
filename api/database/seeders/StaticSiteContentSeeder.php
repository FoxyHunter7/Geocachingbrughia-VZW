<?php

namespace Database\Seeders;

use App\Models\staticSiteContent;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class StaticSiteContentSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('static_site_contents')->truncate();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/staticSiteContents.csv'), ';');

        $model = new staticSiteContent();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
