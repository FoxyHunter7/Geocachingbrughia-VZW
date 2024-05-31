<?php

namespace Database\Seeders;

use App\Models\staticSiteContent;
use Illuminate\Database\Seeder;

class StaticSiteContentSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/staticSiteContents.csv'), ';');

        $model = new staticSiteContent();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
